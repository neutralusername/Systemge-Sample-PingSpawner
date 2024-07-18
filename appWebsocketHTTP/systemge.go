package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Helpers"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Tcp"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: Tcp.NewEndpoint("127.0.0.1:60000", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.PINGPONG: func(node *Node.Node, message *Message.Message) (string, error) {
			println(node.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return "pong", nil
		},
	}
}
