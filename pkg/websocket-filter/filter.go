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
	"log"

	proto "github.com/golang/protobuf/proto"
	ws "github.com/illacloud/builder-backend/internal/websocket"
)

func Run(hub *ws.Hub) {
	for {
		select {
		// handle register event
		case client := <-hub.Register:
			hub.Clients[client.ID] = client
		case client := <-hub.RegisterBinary:
			hub.BinaryClients[client.ID] = client
		// handle unregister events
		case client := <-hub.Unregister:
			if _, ok := hub.Clients[client.ID]; ok {
				delete(hub.Clients, client.ID)
				close(client.Send)
			}
		// handle all hub broadcast events
		case message := <-hub.Broadcast:
			for _, client := range hub.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(hub.Clients, client.ID)
				}
			}
		// handle client on message event
		case message := <-hub.OnTextMessage:
			SignalFilter(hub, message)
		case message := <-hub.OnBinaryMessage:
			BinarySignalFilter(hub, message)
		}

	}
}

func SignalFilter(hub *ws.Hub, message *ws.Message) error {
	switch message.Signal {
	case ws.SIGNAL_PING:
		go SignalPing(hub, message)
	case ws.SIGNAL_ENTER:
		go SignalEnter(hub, message)
	case ws.SIGNAL_LEAVE:
		go SignalLeave(hub, message)
	case ws.SIGNAL_CREATE_STATE:
		go SignalCreateState(hub, message)
	case ws.SIGNAL_DELETE_STATE:
		go SignalDeleteState(hub, message)
	case ws.SIGNAL_UPDATE_STATE:
		go SignalUpdateState(hub, message)
	case ws.SIGNAL_MOVE_STATE:
		go SignalMoveState(hub, message)
	case ws.SIGNAL_CREATE_OR_UPDATE_STATE:
		go SignalCreateOrUpdateState(hub, message)
	case ws.SIGNAL_BROADCAST_ONLY:
		go SignalBroadcastOnly(hub, message)
	case ws.SIGNAL_PUT_STATE:
		go SignalPutState(hub, message)
	case ws.SIGNAL_GLOBAL_BROADCAST_ONLY:
		go SignalGlobalBroadcastOnly(hub, message)
	case ws.SIGNAL_COOPERATE_ATTACH:
		go SignalCooperateAttach(hub, message)
	case ws.SIGNAL_COOPERATE_DISATTACH:
		go SignalCooperateDisattach(hub, message)
	case ws.SIGNAL_SUPER_POWER:
		go SignalEcho(hub, message)
	default:
		return nil
	}
	return nil
}

func BinarySignalFilter(hub *ws.Hub, message []byte) error {
	binaryMessageType, errInGetMessageType := ws.GetBinaryMessageType(message)
	if errInGetMessageType != nil {
		return errInGetMessageType
	}

	switch binaryMessageType {
	case ws.BINARY_MESSAGE_TYPE_MOVING:
		// decode binary message
		movingMessageBin := &ws.MovingMessageBin{}
		if errInParse := proto.Unmarshal(message, movingMessageBin); errInParse != nil {
			log.Printf("[BinarySignalFilter] Failed to parse message MovingMessageBin: ", errInParse)
			return errInParse
		}

		// process message
		MovingMessageFilter(hub, movingMessageBin)

	}
	return nil
}

func MovingMessageFilter(hub *ws.Hub, message *ws.MovingMessageBin) error {
	switch message.Signal {
	case ws.SIGNAL_MOVE_STATE:
		return SignalMoveStateBinary(hub, message)
	case ws.SIGNAL_MOVE_CURSOR:
		return SignalMoveCursorBinary(hub, message)
	default:
		return nil
	}
}

func OptionFilter(hub *ws.Hub, client *ws.Client, message *ws.Message) error {
	return nil
}
