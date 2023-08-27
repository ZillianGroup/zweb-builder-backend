package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/illacloud/builder-backend/src/model"
	"github.com/illacloud/builder-backend/src/request"
	"github.com/illacloud/builder-backend/src/utils/accesscontrol"
	"github.com/illacloud/builder-backend/src/utils/datacontrol"

	"github.com/gin-gonic/gin"
)

func (impl PublicActionRestHandlerImpl) RunPublicAction(c *gin.Context) {
	// fetch needed param
	teamIdentifier, errInGetTeamIdentifier := GetStringParamFromRequest(c, PARAM_TEAM_IDENTIFIER)
	publicActionID, errInGetPublicActionID := controller.GetMagicIntParamFromRequest(c, PARAM_ACTION_ID)
	appID, errInGetAppID := controller.GetMagicIntParamFromRequest(c, PARAM_APP_ID)
	if errInGetTeamIdentifier != nil || errInGetPublicActionID != nil || errInGetAppID != nil {
		return
	}

	// get team id by team teamIdentifier
	team, errInGetTeamInfo := datacontrol.GetTeamInfoByIdentifier(teamIdentifier)
	if errInGetTeamInfo != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get target team by identifier error: "+errInGetTeamInfo.Error())
		return
	}
	teamID := team.GetID()

	// validate
	controller.AttributeGroup.Init()
	controller.AttributeGroup.SetTeamID(teamID)
	controller.AttributeGroup.SetUserAuthToken(accesscontrol.ANONYMOUS_AUTH_TOKEN)
	controller.AttributeGroup.SetUnitType(accesscontrol.UNIT_TYPE_ACTION)
	controller.AttributeGroup.SetUnitID(accesscontrol.DEFAULT_UNIT_ID)
	canManage, errInCheckAttr := controller.AttributeGroup.CanManage(accesscontrol.ACTION_MANAGE_RUN_ACTION)
	if errInCheckAttr != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
		return
	}
	if !canManage {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// check if action is public action
	action, errInRetrieveAction := controller.Storage.ActionStorage.RetrieveActionsByTeamIDActionIDAndVersion(teamID, appID, app.ExportMainlineVersion())
	if errInRetrieveAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_ACTION, "get action failed: "+errInRetrieveAction.Error())
		return
	}
	if !action.IsPublic() {
		controller.FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this action.")
		return
	}

	// set resource timing header
	// @see:
	// [Timing-Allow-Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Timing-Allow-Origin)
	// [Resource_timing](https://developer.mozilla.org/en-US/docs/Web/API/Performance_API/Resource_timing)
	c.Header("Timing-Allow-Origin", "*")

	// execute
	runActionRequest := request.NewRunActionRequest()
	if err := json.NewDecoder(c.Request.Body).Decode(&runActionRequest); err != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error"+err.Error())
		return
	}

	// get action mapped app
	app, errInRetrieveApp := controller.Storage.AppStorage.RetrieveAppByTeamIDAndAppID(teamID, appID)
	if errInRetrieveApp != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_APP, "get app failed: "+errInRetrieveApp.Error())
		return
	}

	// get action
	action, errInRetrieveAction := controller.Storage.ActionStorage.RetrieveActionsByTeamIDActionIDAndVersion(teamID, appID, app.ExportMainlineVersion())
	if errInRetrieveAction != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_ACTION, "get action failed: "+errInRetrieveAction.Error())
		return
	}

	// update action data with run action reqeust
	action.UpdateWithRunActionRequest(runActionRequest, userID)

	// assembly action
	actionFactory := model.NewActionFactoryByAction(action)
	actionAssemblyLine, errInBuild := actionFactory.Build()
	if errInBuild == nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate action type error: "+errInBuild.Error())
		return
	}

	// get resource
	resource := model.NewResource()
	if !action.IsVirtualAction() {
		// process normal resource action
		var errInRetrieveResource error
		resource, errInRetrieveResource = controller.Storage.ResourceStorage.RetrieveByTeamIDAndResourceID(teamID, action.ExportResourceID())
		if errInRetrieveResource != nil {
			controller.FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_RESOURCE, "get resource failed: "+errInRetrieveResource.Error())
			return
		}
		// resource option validate only happend in create or update phrase
	} else {
		// process virtual resource action
		action.AppendRuntimeInfoForVirtualResource(userAuthToken)
	}

	// check action template
	_, errInValidate := actionAssemblyLine.ValidateActionTemplate(action.ExportTemplateInMap())
	if errInValidate != nil {
		controller.FeedbackBadRequest(c, ERROR_FLAG_VALIDATE_REQUEST_BODY_FAILED, "validate action template error: "+errInValidate.Error())
		return errInValidate
	}

	// run
	actionRunResult, errInRunAction := actionAssemblyLine.Run(resource.Options, action.Template)
	if errInRunAction != nil {
		if strings.HasPrefix(errInRunAction.Error(), "Error 1064:") {
			lineNumber, _ := strconv.Atoi(errInRunAction.Error()[len(errInRunAction.Error())-1:])
			message := ""
			regexp, _ := regexp.Compile(`to use`)
			match := regexp.FindStringIndex(errInRunAction.Error())
			if len(match) == 2 {
				message = errInRunAction.Error()[match[1]:]
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode":    400,
				"errorFlag":    ERROR_FLAG_EXECUTE_ACTION_FAILED,
				"errorMessage": errors.New("SQL syntax error").Error(),
				"errorData": map[string]interface{}{
					"lineNumber": lineNumber,
					"message":    "SQL syntax error" + message,
				},
			})
			return
		}
		controller.FeedbackBadRequest(c, ERROR_FLAG_EXECUTE_ACTION_FAILED, "run action error: "+errInRunAction.Error())
		return
	}

	// feedback
	c.JSON(http.StatusOK, actionRunResult)
}
