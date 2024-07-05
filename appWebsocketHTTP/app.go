package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Node"
)

type AppWebsocketHTTP struct {
}

func New() *AppWebsocketHTTP {
	return &AppWebsocketHTTP{}
}

func (app *AppWebsocketHTTP) OnStart(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) GetApplicationConfig() Config.Application {
	return Config.Application{
		HandleMessagesSequentially: false,
	}
}
