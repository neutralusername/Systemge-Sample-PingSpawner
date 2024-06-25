package appPing

import (
	"Systemge/Client"
	"Systemge/Message"
	"SystemgeSamplePingSpawner/topics"
)

type App struct {
	id string
}

func New(id string) Client.Application {
	app := &App{
		id: id,
	}
	return app
}

func (app *App) OnStart(client *Client.Client) error {
	response, err := client.SyncMessage(topics.PING, client.GetName(), "ping")
	if err != nil {
		panic(err)
	}
	println(client.GetName() + " received \"" + response.GetPayload() + "\" from: " + response.GetOrigin())
	return nil
}

func (app *App) OnStop(client *Client.Client) error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		app.id: func(client *Client.Client, message *Message.Message) error {
			println(client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return nil
		},
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{}
}

func (app *App) GetCustomCommandHandlers() map[string]Client.CustomCommandHandler {
	return map[string]Client.CustomCommandHandler{}
}
