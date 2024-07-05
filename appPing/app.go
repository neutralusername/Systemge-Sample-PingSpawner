package appPing

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
)

type App struct {
	id string
}

func New(id string) Node.Application {
	app := &App{
		id: id,
	}
	return app
}

func (app *App) OnStart(node *Node.Node) error {
	println(node.GetName() + " started")
	response, err := node.SyncMessage(topics.PING, node.GetName(), "ping")
	if err != nil {
		panic(err)
	}
	println(node.GetName() + " received \"" + response.GetPayload() + "\" from: " + response.GetOrigin())
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	println(node.GetName() + " ended")
	return nil
}

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

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{}
}

func (app *App) GetApplicationConfig() Config.Application {
	return Config.Application{
		HandleMessagesSequentially: false,
	}
}
