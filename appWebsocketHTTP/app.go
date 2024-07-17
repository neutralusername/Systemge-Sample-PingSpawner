package appWebsocketHTTP

import "Systemge/Node"

type AppWebsocketHTTP struct {
}

func New() *AppWebsocketHTTP {
	return &AppWebsocketHTTP{}
}

func (app *AppWebsocketHTTP) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{}
}
