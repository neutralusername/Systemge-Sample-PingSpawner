package appPing

import (
	"Systemge/Node"
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

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{}
}
