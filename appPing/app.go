package appPing

import (
	"Systemge/Message"
	"Systemge/Node"
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
	response, err := node.SyncMessage(topics.PING, node.GetName(), "ping")
	if err != nil {
		panic(err)
	}
	println(node.GetName() + " received \"" + response.GetPayload() + "\" from: " + response.GetOrigin())
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
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

func (app *App) GetApplicationConfig() Node.ApplicationConfig {
	return Node.ApplicationConfig{
		ResolverAddress:            "127.0.0.1:60000",
		ResolverNameIndication:     "127.0.0.1",
		ResolverTLSCert:            Utilities.GetFileContent("MyCertificate.crt"),
		HandleMessagesSequentially: false,
	}
}
