package appPing

import (
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

func (app *App) OnStart(client *Node.Node) error {
	response, err := client.SyncMessage(topics.PING, client.GetName(), "ping")
	if err != nil {
		panic(err)
	}
	println(client.GetName() + " received \"" + response.GetPayload() + "\" from: " + response.GetOrigin())
	return nil
}

func (app *App) OnStop(client *Node.Node) error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		app.id: func(client *Node.Node, message *Message.Message) error {
			println(client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
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
