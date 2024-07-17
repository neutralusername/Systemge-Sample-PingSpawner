package appPing

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		app.id: func(node *Node.Node, message *Message.Message) error {
			println(node.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return nil
		},
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}

func (app *App) GetSystemgeConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,
	}
}
