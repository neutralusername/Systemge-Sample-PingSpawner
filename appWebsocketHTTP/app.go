package appWebsocketHTTP

import "Systemge/Client"

type AppWebsocketHTTP struct {
}

func New() Client.WebsocketHTTPApplication {
	return &AppWebsocketHTTP{}
}

func (app *AppWebsocketHTTP) OnStart(client *Client.Client) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(client *Client.Client) error {
	return nil
}
