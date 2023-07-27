package filter

import (
	"github.com/illacloud/builder-backend/internal/repository"
	"github.com/illacloud/builder-backend/internal/util/config"
	ws "github.com/illacloud/builder-backend/internal/websocket"
)

func RecordModifyHistory(hub *ws.Hub, message *ws.Message, displayNames []string) error {
	// check deploy env
	conf := config.GetInstance()
	if !conf.IsCloudMode() {
		return nil
	}
	// go
	currentClient, _ := hub.Clients[message.ClientID]
	teamID := currentClient.TeamID
	appID := currentClient.APPID
	userID := currentClient.MappedUserID

	// get current edit version app snapshot
	appSnapshot, errInGetSnapshot := hub.AppSnapshotRepositoryImpl.RetrieveEditVersion(teamID, appID)
	if errInGetSnapshot != nil {
		currentClient.Feedback(message, ws.ERROR_CREATE_SNAPSHOT_MIDIFY_HISTORY_FAILED, errInGetSnapshot)
		return errInGetSnapshot
	}

	// new modify history
	for _, displayName := range displayNames {
		broadcastType := ""
		var broadcastPayload interface{}
		if message.Broadcast != nil {
			broadcastType = message.Broadcast.Type
			broadcastPayload = message.Broadcast.Payload
		}
		modifyHistoryRecord := repository.NewAppModifyHistory(message.Signal, message.Target, displayName, broadcastType, broadcastPayload, userID)
		appSnapshot.PushModifyHistory(modifyHistoryRecord)
	}

	// update app snapshot
	errInUpdateSnapshot := hub.AppSnapshotRepositoryImpl.UpdateWholeSnapshot(appSnapshot)
	if errInUpdateSnapshot != nil {
		currentClient.Feedback(message, ws.ERROR_UPDATE_SNAPSHOT_MIDIFY_HISTORY_FAILED, errInUpdateSnapshot)
		return errInUpdateSnapshot
	}

	// check if app snapshot need archive
	if !appSnapshot.DoesActiveSnapshotNeedArchive() {
		return nil
	}

	// ok, archive app snapshot
	TakeSnapshot(hub, message)

	return nil
}