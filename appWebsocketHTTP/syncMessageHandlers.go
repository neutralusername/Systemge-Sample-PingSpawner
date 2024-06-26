package appWebsocketHTTP

import (
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.PING: func(client *Node.Node, message *Message.Message) (string, error) {
			println(client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}
