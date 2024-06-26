package appWebsocketHTTP

import "Systemge/Node"

type AppWebsocketHTTP struct {
}

func New() Node.WebsocketHTTPApplication {
	return &AppWebsocketHTTP{}
}

func (app *AppWebsocketHTTP) OnStart(client *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(client *Node.Node) error {
	return nil
}
