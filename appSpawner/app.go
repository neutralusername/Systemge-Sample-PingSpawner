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

	spawnedClients map[string]*Client.Client
	mutex          sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client:         client,
		spawnedClients: make(map[string]*Client.Client),
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
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	client := app.spawnedClients[id]
	if client == nil {
		return "", Utilities.NewError("Client "+id+" does not exist", nil)
	}
	err := client.Stop()
	if err != nil {
		return "", Utilities.NewError("Error stopping client "+id, err)
	}
	delete(app.spawnedClients, id)
	err = app.client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error removing async topic \""+id+"\"", err).Error())
	}
	err = app.client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error unregistering topic \""+id+"\"", err).Error())
	}
	//println("ended ping client " + id)
	return "", nil
}

func (app *App) New(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	if _, ok := app.spawnedClients[id]; ok {
		return "", Utilities.NewError("Client "+id+" already exists", nil)
	}
	pingClient := Client.New("client"+id, app.client.GetResolverResolution(), app.client.GetLogger())
	pingApp, err := appPing.New(pingClient, []string{id})
	if err != nil {
		return "", Utilities.NewError("Error creating ping app "+id, err)
	}
	pingClient.SetApplication(pingApp)
	err = app.client.AddAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		return "", Utilities.NewError("Error adding async topic \""+id+"\"", err)
	}
	err = app.client.AddResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), "brokerPing", id)
	if err != nil {
		err = app.client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error removing async topic \""+id+"\"", err).Error())
		}
		return "", Utilities.NewError("Error registering topic", err)
	}
	err = pingClient.Start()
	if err != nil {
		err = app.client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error removing async topic \""+id+"\"", err).Error())
		}
		err = app.client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error unregistering topic \""+id+"\"", err).Error())
		}
		return "", Utilities.NewError("Error starting client", err)
	}
	//println("created ping client " + id)
	app.spawnedClients[id] = pingClient
	return id, nil
}
