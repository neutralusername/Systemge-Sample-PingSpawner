package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{}
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketClient.Client) {
	_, err := app.client.SyncMessage(topics.NEW, websocketClient.GetId(), websocketClient.GetId())
	if err != nil {
		panic(Utilities.NewError("Error sending sync message", err))
	}
	err = app.client.AsyncMessage(websocketClient.GetId(), websocketClient.GetId(), "ping")
	if err != nil {
		panic(Utilities.NewError("Error sending async message", err))
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(websocketClient *WebsocketClient.Client) {
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), websocketClient.GetId())
	if err != nil {
		//windows seems to have issues with the sync token generation.. sometimes it will generate two similar tokens in sequence. i assume the system time is not accurate enough for very fast token generation
		panic(Utilities.NewError("Error sending sync message", err))
	}
}
