package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetSystemgeConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.PING: func(node *Node.Node, message *Message.Message) (string, error) {
			println(node.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}
