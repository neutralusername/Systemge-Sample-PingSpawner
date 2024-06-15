package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSamplePingSpawner/topics"
)

type WebsocketApp struct {
	client        *Client.Client
	clientPingIds map[string]string
}

func New(messageBrokerClient *Client.Client, args []string) (Application.WebsocketApplication, error) {
	return &WebsocketApp{
		client:        messageBrokerClient,
		clientPingIds: make(map[string]string),
	}, nil
}

func (app *WebsocketApp) OnStart() error {
	return nil
}

func (app *WebsocketApp) OnStop() error {
	return nil
}

func (app *WebsocketApp) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{}
}

func (app *WebsocketApp) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.PING: func(message *Message.Message) (string, error) {
			println(app.client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}

func (app *WebsocketApp) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *WebsocketApp) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{}
}

func (app *WebsocketApp) OnConnectHandler(connection *WebsocketClient.Client) {
	response, err := app.client.SyncMessage(topics.NEW, connection.GetId(), "")
	if err != nil {
		panic(Utilities.NewError("Error sending sync message", err))
	}
	app.clientPingIds[connection.GetId()] = response.GetPayload()
	err = app.client.AsyncMessage("ping_"+response.GetPayload(), connection.GetId(), "ping")
	if err != nil {
		panic(Utilities.NewError("Error sending async message", err))
	}
}

func (app *WebsocketApp) OnDisconnectHandler(connection *WebsocketClient.Client) {
	pingId := app.clientPingIds[connection.GetId()]
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), pingId)
	if err != nil {
		panic(err)
	}
}
