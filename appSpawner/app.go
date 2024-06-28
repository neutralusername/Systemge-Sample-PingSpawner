package appSpawner

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Resolution"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/topics"
	"sync"
)

type App struct {
	spawnedNodes map[string]*Node.Node
	mutex        sync.Mutex
}

func New() Node.Application {
	app := &App{
		spawnedNodes: make(map[string]*Node.Node),
	}
	return app
}

func (app *App) OnStart(node *Node.Node) error {
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
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

func (app *App) End(node *Node.Node, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	spawnedNode := app.spawnedNodes[id]
	if spawnedNode == nil {
		return "", Error.New("Node "+id+" does not exist", nil)
	}
	err := spawnedNode.Stop()
	if err != nil {
		return "", Error.New("Error stopping node "+id, err)
	}
	delete(app.spawnedNodes, id)
	err = node.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		node.GetLogger().Log(Error.New("Error removing async topic \""+id+"\"", err).Error())
	}
	err = node.RemoveResolverTopicRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		node.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
	}
	println("ended pingNode " + id)
	return "", nil
}

func (app *App) New(node *Node.Node, message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	if _, ok := app.spawnedNodes[id]; ok {
		return "", Error.New("Node "+id+" already exists", nil)
	}
	pingNode := Node.New(Config.Node{
		Name:       id,
		LoggerPath: "error.log",
	}, appPing.New(id), nil, nil)

	err := node.AddAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		return "", Error.New("Error adding async topic \""+id+"\"", err)
	}
	err = node.AddResolverTopicRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), *Resolution.New("brokerPing", "127.0.0.1:60007", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt")), id)
	if err != nil {
		errRemoveAsync := node.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if errRemoveAsync != nil {
			node.GetLogger().Log(Error.New("Error removing async topic \""+id+"\"", errRemoveAsync).Error())
		}
		return "", Error.New("Error registering topic", err)
	}
	err = pingNode.Start()
	if err != nil {
		errRemoveAsync := node.RemoveAsyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if errRemoveAsync != nil {
			node.GetLogger().Log(Error.New("Error removing async topic \""+id+"\"", errRemoveAsync).Error())
		}
		errRemoveResolver := node.RemoveResolverTopicRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if errRemoveResolver != nil {
			node.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", errRemoveResolver).Error())
		}
		return "", Error.New("Error starting node", err)
	}
	println("created pingNode " + id)
	app.spawnedNodes[id] = pingNode
	return id, nil
}

func (app *App) GetApplicationConfig() Config.Application {
	return Config.Application{
		ResolverAddress:            "127.0.0.1:60000",
		ResolverNameIndication:     "127.0.0.1",
		ResolverTLSCert:            Utilities.GetFileContent("MyCertificate.crt"),
		HandleMessagesSequentially: false,
	}
}
