// Copyright 2022 The ZWEB Authors.
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

package filter

import (
	"errors"

	"github.com/zilliangroup/zweb-builder-backend/src/model"
	"github.com/zilliangroup/zweb-builder-backend/src/websocket"

	"github.com/zilliangroup/zweb-builder-backend/src/utils/builderoperation"
)

func SignalCreateState(hub *websocket.Hub, message *websocket.Message) error {
	// init global param
	currentClient, errInGetClient := hub.GetClientByID(message.ClientID)
	if errInGetClient != nil {
		return errInGetClient
	}
	stateType := model.STATE_TYPE_INVALIED
	teamID := currentClient.TeamID
	appID := currentClient.APPID
	userID := currentClient.MappedUserID

	// rewrite message broadcast
	message.RewriteBroadcast()

	// new app
	app, errInRetrieveApp := hub.Storage.AppStorage.RetrieveAppByTeamIDAndAppID(teamID, appID)
	if errInRetrieveApp != nil {
		currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInRetrieveApp)
		return errInRetrieveApp
	}
	app.Modify(userID)

	// modified displayNames
	displayNames := make([]string, 0)

	// target switch
	switch message.Target {
	case builderoperation.TARGET_NOTNING:
		return nil
	case builderoperation.TARGET_COMPONENTS:
		stateType = model.TREE_STATE_TYPE_COMPONENTS
		for _, v := range message.Payload {
			// init current tree state node
			currentTreeStateNode, errInInitCurrentNode := model.NewTreeStateByWebsocketMessage(app, stateType, v)
			if errInInitCurrentNode != nil {
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInInitCurrentNode)
				return errInInitCurrentNode
			}
			// current state did not in database, create
			componentTree := model.ConstructComponentNodeByMap(v)
			errInCreateComponentTree := hub.Storage.TreeStateStorage.CreateComponentTree(app, model.TREE_STATE_SUMMIT_ID, componentTree)
			if errInCreateComponentTree != nil {
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInCreateComponentTree)
				return errInCreateComponentTree
			}
			displayNames = append(displayNames, currentTreeStateNode.ExportName())
		}
	case builderoperation.TARGET_DEPENDENCIES:
		stateType = model.KV_STATE_TYPE_DEPENDENCIES
		// create k-v state
		for _, v := range message.Payload {
			subv, ok := v.(map[string]interface{})
			if !ok {
				errInResolveSubv := errors.New("K-V State reflect failed, please check your input.")
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInResolveSubv)
				return errInResolveSubv
			}
			for key, depState := range subv {
				// init current kvState node
				currentKVStateNode, errInNewKVState := model.NewKVStateByWebsocketMessageWithGivenKey(app, stateType, key, depState)
				if errInNewKVState != nil {
					currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInNewKVState)
					return errInNewKVState
				}
				// current state did not in database, create
				errInCreateKVState := hub.Storage.KVStateStorage.Create(currentKVStateNode)
				if errInCreateKVState != nil {
					currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInCreateKVState)
					return errInCreateKVState
				}
				displayNames = append(displayNames, currentKVStateNode.ExportKey())
			}
		}
	case builderoperation.TARGET_DRAG_SHADOW:
		fallthrough

	case builderoperation.TARGET_DOTTED_LINE_SQUARE:
		// fill type
		if message.Target == builderoperation.TARGET_DRAG_SHADOW {
			stateType = model.KV_STATE_TYPE_DRAG_SHADOW
		} else {
			stateType = model.KV_STATE_TYPE_DOTTED_LINE_SQUARE
		}
		// create k-v state
		for _, v := range message.Payload {
			// init current kvState node
			currentKVStateNode, errInNewKVState := model.NewKVStateByWebsocketMessage(app, stateType, v)
			if errInNewKVState != nil {
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInNewKVState)
				return errInNewKVState
			}

			// current state did not in database, create
			errInCreateKVState := hub.Storage.KVStateStorage.Create(currentKVStateNode)
			if errInCreateKVState != nil {
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInCreateKVState)
				return errInCreateKVState
			}
			displayNames = append(displayNames, currentKVStateNode.ExportKey())
		}
	case builderoperation.TARGET_DISPLAY_NAME:
		stateType = model.SET_STATE_TYPE_DISPLAY_NAME
		// create set state
		for _, v := range message.Payload {
			// init current set state node
			currentSetStateNode, errInNewSetState := model.NewSetStateByWebsocketMessage(app, stateType, v)
			if errInNewSetState != nil {
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInNewSetState)
				return errInNewSetState
			}

			// create
			errInCreateSetState := hub.Storage.SetStateStorage.Create(currentSetStateNode)
			if errInCreateSetState != nil {
				currentClient.Feedback(message, websocket.ERROR_CREATE_STATE_FAILED, errInCreateSetState)
				return errInCreateSetState
			}
			displayNames = append(displayNames, currentSetStateNode.ExportValue())
		}
	case builderoperation.TARGET_APPS:
		// serve on HTTP API, this signal only for broadcast
		for _, v := range message.Payload {
			appForExport, errInNewAppForExport := model.NewAppForExportByMap(v)
			if errInNewAppForExport == nil {
				displayNames = append(displayNames, appForExport.ExportName())
			}
		}
	case builderoperation.TARGET_RESOURCE:
		// serve on HTTP API, this signal only for broadcast
		for _, v := range message.Payload {
			resourceForExport, errInNewResourceForExport := model.NewResourceForExportByMap(v)
			if errInNewResourceForExport == nil {
				displayNames = append(displayNames, resourceForExport.ExportName())
			}
		}
	case builderoperation.TARGET_ACTION:
		// serve on HTTP API, this signal only for broadcast
		for _, v := range message.Payload {
			actionForExport, errInNewActionForExport := model.NewActionForExportByMap(v)
			if errInNewActionForExport == nil {
				displayNames = append(displayNames, actionForExport.ExportDisplayName())
			}
		}
	}

	// record app snapshot modify history
	RecordModifyHistory(hub, message, displayNames)

	// the currentClient does not need feedback when operation success

	// change app modify time
	hub.Storage.AppStorage.UpdateWholeApp(app)

	// feedback otherClient
	hub.BroadcastToOtherClients(message, currentClient)
	return nil
}
