package appChess

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
)

type App struct {
	client *Client.Client

	moveTopic string
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	if len(args) != 1 {
		return nil, Utilities.NewError("Expected 1 argument", nil)
	}
	app := &App{
		client: client,

		moveTopic: args[0],
	}
	return app, nil
}

func (app *App) OnStart() error {
	println("chess app started", app.moveTopic)
	return nil
}

func (app *App) OnStop() error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		app.moveTopic: func(message *Message.Message) error {
			return nil
		},
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}
