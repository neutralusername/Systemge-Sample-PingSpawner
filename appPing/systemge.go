package appPing

import (
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		"pong": func(node *Node.Node, message *Message.Message) error {
			println(node.GetName() + " received pong-async")
			return nil
		},
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		"ping": func(node *Node.Node, message *Message.Message) (string, error) {
			println(node.GetName() + " received ping-sync; sending pong-sync")
			return "", nil
		},
	}
}
