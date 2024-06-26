package appWebsocketHTTP

import "Systemge/Node"

type AppWebsocketHTTP struct {
}

func New() Node.WebsocketHTTPApplication {
	return &AppWebsocketHTTP{}
}

func (app *AppWebsocketHTTP) OnStart(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(node *Node.Node) error {
	return nil
}
