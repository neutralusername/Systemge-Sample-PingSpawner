package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topics.PING: func(client *Client.Client, message *Message.Message) (string, error) {
			println(client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}
