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

package filter

import (
	"errors"

	"github.com/illacloud/builder-backend/internal/repository"
	"github.com/illacloud/builder-backend/pkg/app"
	"github.com/illacloud/builder-backend/pkg/state"

	"github.com/illacloud/builder-backend/internal/util/builderoperation"
	ws "github.com/illacloud/builder-backend/internal/websocket"
)

func SignalDeleteState(hub *ws.Hub, message *ws.Message) error {

	// deserialize message
	currentClient, hit := hub.Clients[message.ClientID]
	if !hit {
		return errors.New("[SignalDeleteState] target client(" + message.ClientID.String() + ") does dot exists.")
	}
	stateType := repository.STATE_TYPE_INVALIED
	teamID := currentClient.TeamID
	appDto := app.NewAppDto()
	appDto.ConstructWithID(currentClient.APPID)
	appDto.ConstructWithUpdateBy(currentClient.MappedUserID)
	appDto.SetTeamID(currentClient.TeamID)
	message.RewriteBroadcast()

	// target switch
	switch message.Target {
	case builderoperation.TARGET_NOTNING:
		return nil
	case builderoperation.TARGET_COMPONENTS:
		displayNames := make([]string, 0)
		for _, v := range message.Payload {
			currentNode := state.NewTreeStateDto()
			currentNode.InitUID()
			currentNode.SetTeamID(teamID)
			currentNode.ConstructWithDisplayNameForDelete(v) // set Name
			currentNode.ConstructByApp(appDto)               // set AppRefID
			currentNode.ConstructWithType(repository.TREE_STATE_TYPE_COMPONENTS)

			if err := hub.TreeStateServiceImpl.DeleteTreeStateNodeRecursive(currentNode); err != nil {
				currentClient.Feedback(message, ws.ERROR_DELETE_STATE_FAILED, err)
				return err
			}
			// collect display names
			displayNames = append(displayNames, currentNode.ExportName())
		}
		// record app snapshot modify history
		RecordModifyHistory(hub, message, displayNames)
	case builderoperation.TARGET_DEPENDENCIES:
		// dependency can not delete

	case builderoperation.TARGET_DRAG_SHADOW:
		fallthrough

	case builderoperation.TARGET_DOTTED_LINE_SQUARE:
		// fill type
		if message.Target == builderoperation.TARGET_DRAG_SHADOW {
			stateType = repository.KV_STATE_TYPE_DRAG_SHADOW
		} else {
			stateType = repository.KV_STATE_TYPE_DOTTED_LINE_SQUARE
		}
		// delete k-v state
		for _, v := range message.Payload {
			// fill KVStateDto
			kvStateDto := state.NewKVStateDto()
			kvStateDto.InitUID()
			kvStateDto.SetTeamID(teamID)
			kvStateDto.ConstructWithDisplayNameForDelete(v)
			kvStateDto.ConstructByApp(appDto) // set AppRefID
			kvStateDto.ConstructWithType(stateType)

			if err := hub.KVStateServiceImpl.DeleteKVStateByKey(kvStateDto); err != nil {
				currentClient.Feedback(message, ws.ERROR_DELETE_STATE_FAILED, err)
				return err
			}
		}

	case builderoperation.TARGET_DISPLAY_NAME:
		stateType = repository.SET_STATE_TYPE_DISPLAY_NAME
		displayNames := make([]string, 0)
		// delete set state
		for _, v := range message.Payload {
			// resolve payload
			displayName, err := repository.ResolveDisplayNameByPayload(v)
			if err != nil {
				return err
			}
			// init
			setStateDto := state.NewSetStateDto()
			setStateDto.InitUID()
			setStateDto.SetTeamID(teamID)
			setStateDto.ConstructWithDisplayNameForDelete(v)
			setStateDto.ConstructWithType(stateType)
			setStateDto.ConstructByApp(appDto)
			setStateDto.ConstructWithEditVersion()
			// delete state
			if err := hub.SetStateServiceImpl.DeleteSetStateByValue(setStateDto); err != nil {
				currentClient.Feedback(message, ws.ERROR_CREATE_STATE_FAILED, err)
				return err
			}
			displayNames = append(displayNames, displayName)
		}
	case builderoperation.TARGET_APPS:
		// serve on HTTP API, this signal only for broadcast
		displayNames := make([]string, 0)
		for _, v := range message.Payload {
			appForExport, errInNewAppForExport := repository.NewAppForExportByMap(v)
			if errInNewAppForExport == nil {
				displayNames = append(displayNames, appForExport.ExportName())
			}
		}
		// record app snapshot modify history
		RecordModifyHistory(hub, message, displayNames)
	case builderoperation.TARGET_RESOURCE:
		// serve on HTTP API, this signal only for broadcast
		displayNames := make([]string, 0)
		for _, v := range message.Payload {
			resourceForExport, errInNewResourceForExport := repository.NewResourceForExportByMap(v)
			if errInNewResourceForExport == nil {
				displayNames = append(displayNames, resourceForExport.ExportName())
			}
		}
		// record app snapshot modify history
		RecordModifyHistory(hub, message, displayNames)
	case builderoperation.TARGET_ACTION:
		// serve on HTTP API, this signal only for broadcast
		displayNames := make([]string, 0)
		for _, v := range message.Payload {
			actionForExport, errInNewActionForExport := repository.NewActionForExportByMap(v)
			if errInNewActionForExport == nil {
				displayNames = append(displayNames, actionForExport.ExportDisplayName())
			}
		}
		// record app snapshot modify history
		RecordModifyHistory(hub, message, displayNames)
	}

	// the currentClient does not need feedback when operation success

	// change app modify time
	hub.AppServiceImpl.UpdateAppModifyTime(appDto)

	// feedback otherClient
	hub.BroadcastToOtherClients(message, currentClient)

	return nil
}
