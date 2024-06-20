package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSamplePingSpawner/topics"
)

type AppWebsocket struct {
	client *Client.Client
}

func New(messageBrokerClient *Client.Client, args []string) (Application.WebsocketApplication, error) {
	return &AppWebsocket{
		client: messageBrokerClient,
	}, nil
}

func (app *AppWebsocket) OnStart() error {
	return nil
}

func (app *AppWebsocket) OnStop() error {
	return nil
}

func (app *AppWebsocket) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{}
}

func (app *AppWebsocket) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.PING: func(message *Message.Message) (string, error) {
			println(app.client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}

func (app *AppWebsocket) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *AppWebsocket) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{}
}

func (app *AppWebsocket) OnConnectHandler(connection *WebsocketClient.Client) {
	_, err := app.client.SyncMessage(topics.NEW, connection.GetId(), connection.GetId())
	if err != nil {
		panic(Utilities.NewError("Error sending sync message", err))
	}
	err = app.client.AsyncMessage(connection.GetId(), connection.GetId(), "ping")
	if err != nil {
		panic(Utilities.NewError("Error sending async message", err))
	}
}

func (app *AppWebsocket) OnDisconnectHandler(connection *WebsocketClient.Client) {
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), connection.GetId())
	if err != nil {
		//windows seems to have issues with the sync token generation.. sometimes it will generate two similar tokens in sequence. i assume the system time is not accurate enough for very fast token generation
		panic(err)
	}
}
