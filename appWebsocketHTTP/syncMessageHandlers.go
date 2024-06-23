package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Message"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.PING: func(message *Message.Message) (string, error) {
			//println(app.client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}
