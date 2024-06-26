package appSpawner

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/topics"
	"sync"
)

type App struct {
	spawnedClients map[string]*Node.Node
	mutex          sync.Mutex
}

func New() Node.Application {
	app := &App{
		spawnedClients: make(map[string]*Node.Node),
	}
	return app
}

func (app *App) OnStart(client *Node.Node) error {
	return nil
}

func (app *App) OnStop(client *Node.Node) error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.NEW: app.New,
		topics.END: app.End,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{}
}

func (app *App) End(client *Node.Node, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	spawnedClient := app.spawnedClients[id]
	if spawnedClient == nil {
		return "", Error.New("Node "+id+" does not exist", nil)
	}
	err := spawnedClient.Stop()
	if err != nil {
		return "", Error.New("Error stopping client "+id, err)
	}
	delete(app.spawnedClients, id)
	err = client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		client.GetLogger().Log(Error.New("Error removing async topic \""+id+"\"", err).Error())
	}
	err = client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		client.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
	}
	println("ended ping client " + id)
	return "", nil
}

func (app *App) New(client *Node.Node, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	if _, ok := app.spawnedClients[id]; ok {
		return "", Error.New("Node "+id+" already exists", nil)
	}
	pingClientConfig := &Node.Config{
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
		return "", Error.New("Error adding async topic \""+id+"\"", err)
	}
	err = client.AddResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), "brokerPing", id)
	if err != nil {
		err = client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Error.New("Error removing async topic \""+id+"\"", err).Error())
		}
		return "", Error.New("Error registering topic", err)
	}
	err = pingClient.Start()
	if err != nil {
		err = client.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Error.New("Error removing async topic \""+id+"\"", err).Error())
		}
		err = client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
		}
		return "", Error.New("Error starting client", err)
	}
	println("created ping client " + id)
	app.spawnedClients[id] = pingClient
	return id, nil
}
