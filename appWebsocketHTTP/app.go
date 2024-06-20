package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Client"
)

type AppWebsocketHTTP struct {
	client *Client.Client
}

func New(client *Client.Client, args []string) (Application.CompositeApplicationWebsocketHTTP, error) {
	return &AppWebsocketHTTP{
		client: client,
	}, nil
}

func (app *AppWebsocketHTTP) OnStart() error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop() error {
	return nil
}
