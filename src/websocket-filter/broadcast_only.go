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

	"github.com/zilliangroup/zweb-builder-backend/src/websocket"
)

func SignalBroadcastOnly(hub *websocket.Hub, message *websocket.Message) error {
	// deserialize message
	currentClient, hit := hub.Clients[message.ClientID]
	if !hit {
		return errors.New("[SignalBroadcastOnly] target client(" + message.ClientID.String() + ") does dot exists.")
	}
	message.RewriteBroadcast()

	// feedback otherClient
	hub.BroadcastToOtherClients(message, currentClient)
	return nil
}
