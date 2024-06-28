package appPing

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Resolution"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/topics"
)

type App struct {
	id string
}

func New(id string) Node.Application {
	app := &App{
		id: id,
	}
	return app
}

func (app *App) OnStart(node *Node.Node) error {
	println(node.GetName() + " started")
	response, err := node.SyncMessage(topics.PING, node.GetName(), "ping")
	if err != nil {
		panic(err)
	}
	println(node.GetName() + " received \"" + response.GetPayload() + "\" from: " + response.GetOrigin())
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	println(node.GetName() + " stopped")
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		app.id: func(node *Node.Node, message *Message.Message) error {
			println(node.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return nil
		},
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{}
}

func (app *App) GetApplicationConfig() Config.Application {
	return Config.Application{
		ResolverResolution:         Resolution.New("resolver", "127.0.0.1:60000", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
		HandleMessagesSequentially: false,
	}
}
