package appWebsocketHTTP

import "Systemge/Application"

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{}
}
