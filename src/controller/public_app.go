// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by publicApplicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resthandler

import (
	ac "github.com/illacloud/builder-backend/internal/accesscontrol"
	"github.com/illacloud/builder-backend/internal/auditlogger"
	dc "github.com/illacloud/builder-backend/internal/datacontrol"
	"github.com/illacloud/builder-backend/internal/repository"
	"github.com/illacloud/builder-backend/pkg/app"
	"github.com/illacloud/builder-backend/pkg/state"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PublicAppRestHandler interface {
	GetMegaData(c *gin.Context)
	IsPublicApp(c *gin.Context)
}

type PublicAppRestHandlerImpl struct {
	logger           *zap.SugaredLogger
	appService       app.AppService
	AttributeGroup   *ac.AttributeGroup
	treeStateService state.TreeStateService
}

func NewPublicAppRestHandlerImpl(logger *zap.SugaredLogger, appService app.AppService, attrg *ac.AttributeGroup, treeStateService state.TreeStateService) *PublicAppRestHandlerImpl {
	return &PublicAppRestHandlerImpl{
		logger:           logger,
		appService:       appService,
		AttributeGroup:   attrg,
		treeStateService: treeStateService,
	}
}

func (impl PublicAppRestHandlerImpl) GetMegaData(c *gin.Context) {
	// fetch needed param
	teamIdentifier, errInGetTeamIdentifier := GetStringParamFromRequest(c, PARAM_TEAM_IDENTIFIER)
	publicAppID, errInGetAPPID := GetMagicIntParamFromRequest(c, PARAM_APP_ID)
	version, errInGetVersion := GetIntParamFromRequest(c, PARAM_VERSION)
	if errInGetTeamIdentifier != nil || errInGetAPPID != nil || errInGetVersion != nil {
		return
	}

	// check version, the version must be repository.APP_AUTO_RELEASE_VERSION
	if version != repository.APP_AUTO_RELEASE_VERSION {
		FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you only can access release version of app.")
		return
	}

	// get team id by team teamIdentifier
	team, errInGetTeamInfo := dc.GetTeamInfoByIdentifier(teamIdentifier)
	if errInGetTeamInfo != nil {
		FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_TEAM, "get target team by identifier error: "+errInGetTeamInfo.Error())
		return
	}
	teamID := team.GetID()

	// validate
	impl.AttributeGroup.Init()
	impl.AttributeGroup.SetTeamID(teamID)
	impl.AttributeGroup.SetUserAuthToken(ac.ANONYMOUS_AUTH_TOKEN)
	impl.AttributeGroup.SetUnitType(ac.UNIT_TYPE_APP)
	impl.AttributeGroup.SetUnitID(ac.DEFAULT_UNIT_ID)
	canAccess, errInCheckAttr := impl.AttributeGroup.CanAccess(ac.ACTION_ACCESS_VIEW)
	if errInCheckAttr != nil {
		FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "error in check attribute: "+errInCheckAttr.Error())
		return
	}
	if !canAccess {
		FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this attribute due to access control policy.")
		return
	}

	// check if app is public app
	if !impl.appService.IsPublicApp(teamID, publicAppID) {
		FeedbackBadRequest(c, ERROR_FLAG_ACCESS_DENIED, "you can not access this app.")
		return
	}

	// Fetch Mega data via `publicApp` and `version`
	res, err := impl.appService.GetMegaData(teamID, publicAppID, version)
	if err != nil {
		if err.Error() == "content not found" {
			FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_APP, "get publicApp mega data error: "+err.Error())
			return
		}
		FeedbackBadRequest(c, ERROR_FLAG_CAN_NOT_GET_APP, "get publicApp mega data error: "+err.Error())
		return
	}

	// audit log
	auditLogger := auditlogger.GetInstance()
	auditLogger.Log(&auditlogger.LogInfo{
		EventType: auditlogger.AUDIT_LOG_VIEW_APP,
		TeamID:    teamID,
		UserID:    -1,
		IP:        c.ClientIP(),
		AppID:     publicAppID,
		AppName:   res.AppInfo.Name,
	})

	// feedback
	FeedbackOK(c, res)
	return
}

func (impl PublicAppRestHandlerImpl) IsPublicApp(c *gin.Context) {
	// fetch needed param
	teamIdentifier, errInGetTeamIdentifier := GetStringParamFromRequest(c, PARAM_TEAM_IDENTIFIER)
	publicAppID, errInGetAPPID := GetMagicIntParamFromRequest(c, PARAM_APP_ID)
	if errInGetTeamIdentifier != nil || errInGetAPPID != nil {
		return
	}

	// get team id by team teamIdentifier
	team, errInGetTeamInfo := dc.GetTeamInfoByIdentifier(teamIdentifier)
	if errInGetTeamInfo != nil {
		FeedbackOK(c, repository.NewIsPublicAppResponse(false))
		return
	}
	teamID := team.GetID()

	// check if app is public app
	if !impl.appService.IsPublicApp(teamID, publicAppID) {
		FeedbackOK(c, repository.NewIsPublicAppResponse(false))
		return
	}

	// feedback
	FeedbackOK(c, repository.NewIsPublicAppResponse(true))
	return
}