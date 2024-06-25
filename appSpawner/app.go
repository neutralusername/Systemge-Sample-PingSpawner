package appSpawner

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Module"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/topics"
	"sync"
)

type App struct {
	spawnedClients map[string]*Client.Client
	mutex          sync.Mutex
}

func New() Client.Application {
	app := &App{
		spawnedClients: make(map[string]*Client.Client),
	}
	return app
}

func (app *App) OnStart(client *Client.Client) error {
	return nil
}

func (app *App) OnStop(client *Client.Client) error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topics.NEW: app.New,
		topics.END: app.End,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Client.CustomCommandHandler {
	return map[string]Client.CustomCommandHandler{}
}

func (app *App) End(client *Client.Client, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	spawnedClient := app.spawnedClients[id]
	if spawnedClient == nil {
		return "", Utilities.NewError("Client "+id+" does not exist", nil)
	}
	err := spawnedClient.Stop()
	if err != nil {
		return "", Utilities.NewError("Error stopping client "+id, err)
	}
	delete(app.spawnedClients, id)
	err = client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		client.GetLogger().Log(Utilities.NewError("Error removing async topic \""+id+"\"", err).Error())
	}
	err = client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		client.GetLogger().Log(Utilities.NewError("Error unregistering topic \""+id+"\"", err).Error())
	}
	println("ended ping client " + id)
	return "", nil
}

func (app *App) New(client *Client.Client, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	if _, ok := app.spawnedClients[id]; ok {
		return "", Utilities.NewError("Client "+id+" already exists", nil)
	}
	pingClientConfig := &Client.Config{
		Name:                   id,
		LoggerPath:             "error.log",
		ResolverAddress:        client.GetResolverAddress(),
		ResolverNameIndication: client.GetResolverNameIndication(),
		ResolverTLSCert:        client.GetResolverTLSCert(),
	}
	pingApp := appPing.New(id)
	pingClient := Module.NewClient(pingClientConfig, pingApp, nil, nil)

	err := client.AddAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		return "", Utilities.NewError("Error adding async topic \""+id+"\"", err)
	}
	err = client.AddResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), "brokerPing", id)
	if err != nil {
		err = client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Utilities.NewError("Error removing async topic \""+id+"\"", err).Error())
		}
		return "", Utilities.NewError("Error registering topic", err)
	}
	err = pingClient.Start()
	if err != nil {
		err = client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Utilities.NewError("Error removing async topic \""+id+"\"", err).Error())
		}
		err = client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Utilities.NewError("Error unregistering topic \""+id+"\"", err).Error())
		}
		return "", Utilities.NewError("Error starting client", err)
	}
	println("created ping client " + id)
	app.spawnedClients[id] = pingClient
	return id, nil
}
