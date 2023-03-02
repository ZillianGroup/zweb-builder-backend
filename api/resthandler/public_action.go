// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resthandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	ac "github.com/illacloud/illa-builder-backend/internal/accesscontrol"
	dc "github.com/illacloud/illa-builder-backend/internal/datacontrol"
	"github.com/illacloud/illa-builder-backend/pkg/action"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PublicActionRestHandler interface {
	RunAction(c *gin.Context)
}

type PublicActionRestHandlerImpl struct {
	logger         *zap.SugaredLogger
	actionService  action.ActionService
	AttributeGroup *ac.AttributeGroup
}

func NewPublicActionRestHandlerImpl(logger *zap.SugaredLogger, actionService action.ActionService, attrg *ac.AttributeGroup) *PublicActionRestHandlerImpl {
	return &PublicActionRestHandlerImpl{
		logger:         logger,
		actionService:  actionService,
		AttributeGroup: attrg,
	}
}

func (impl PublicActionRestHandlerImpl) RunAction(c *gin.Context) {
	// fetch needed param
	teamIdentifier, errInGetTeamIdentifier := GetStringParamFromRequest(c, PARAM_TEAM_IDENTIFIER)
	publicActionID, errInGetPublicActionID := GetMagicIntParamFromRequest(c, PARAM_ACTION_ID)
	if errInGetTeamIdentifier != nil || errInGetPublicActionID != nil {
		return
	}

	// get team id by team teamIdentifier
	team, errInGetTeamInfo := dc.GetTeamInfoByIdentifier(teamIdentifier)
	if errInGetTeamInfo != nil {
		FeedbackInternalServerError(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get target team by identifier error: "+errInGetTeamInfo.Error())
		return
	}
	teamID := team.GetID()

	// check if action is public action
	if !impl.actionService.IsPublicAction(teamID, publicActionID) {
		FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this action.")
		return
	}

	// execute
	c.Header("Timing-Allow-Origin", "*")
	var actForExport action.ActionDtoForExport
	if err := json.NewDecoder(c.Request.Body).Decode(&actForExport); err != nil {
		FeedbackBadRequest(c, ERROR_FLAG_PARSE_REQUEST_BODY_FAILED, "parse request body error"+err.Error())
		return
	}
	act := actForExport.ExportActionDto()
	res, err := impl.actionService.RunAction(teamID, act)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1064:") {
			lineNumber, _ := strconv.Atoi(err.Error()[len(err.Error())-1:])
			message := ""
			regexp, _ := regexp.Compile(`to use`)
			match := regexp.FindStringIndex(err.Error())
			if len(match) == 2 {
				message = err.Error()[match[1]:]
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
		FeedbackBadRequest(c, ERROR_FLAG_EXECUTE_ACTION_FAILED, "run action error: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}
