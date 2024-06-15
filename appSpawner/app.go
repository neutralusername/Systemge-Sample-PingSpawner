package appSpawner

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/topics"
	"sync"
)

type App struct {
	client *Client.Client

	idCounter int
	mutex     sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client: client,

		idCounter: 0,
	}
	return app, nil
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.NEW: app.New,
		topics.END: app.End,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *App) End(message *Message.Message) (string, error) {
	return "", nil
}

func (app *App) New(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := Utilities.IntToString(app.idCounter)
	app.idCounter++
	pingTopic := "ping_" + id
	pingClient := Client.New("clientPing"+id, app.client.GetTopicResolutionServerAddress(), app.client.GetLogger(), nil)
	pingApp, err := appPing.New(pingClient, []string{pingTopic})
	if err != nil {
		return "", Utilities.NewError("Error creating ping app "+id, err)
	}
	pingClient.SetApplication(pingApp)
	brokerNetConn, err := Utilities.TlsDial("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return "", Utilities.NewError("Error dialing ping broker", err)
	}
	_, err = Utilities.TcpExchange(brokerNetConn, Message.NewAsync("addAsyncTopic", app.client.GetName(), pingTopic), 5000)
	if err != nil {
		return "", Utilities.NewError("Error exchanging messages with chess broker", err)
	}
	resolverNetConn, err := Utilities.TlsDial("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return "", Utilities.NewError("Error dialing topic resolution server", err)
	}
	_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("registerTopics", app.client.GetName(), "brokerPing "+pingTopic), 5000)
	if err != nil {
		return "", Utilities.NewError("Error exchanging messages with topic resolution server", err)
	}
	err = pingClient.Start()
	if err != nil {
		return "", Utilities.NewError("Error starting chess client", err)
	}

	return "", nil
}
