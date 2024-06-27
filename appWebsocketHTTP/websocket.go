package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	_, err := node.SyncMessage(topics.NEW, websocketClient.GetId(), websocketClient.GetId())
	if err != nil {
		panic(Error.New("Error sending sync message", err))
	}
	err = node.AsyncMessage(websocketClient.GetId(), websocketClient.GetId(), "ping")
	if err != nil {
		panic(Error.New("Error sending async message", err))
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	_, err := node.SyncMessage(topics.END, node.GetName(), websocketClient.GetId())
	if err != nil {
		//windows seems to have issues with the sync token generation.. sometimes it will generate two similar tokens in sequence. i assume the system time is not accurate enough for very fast token generation
		panic(Error.New("Error sending sync message", err))
	}
	node.RemoveTopicResolution(websocketClient.GetId())
}

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() Node.WebsocketComponentConfig {
	return Node.WebsocketComponentConfig{
		Pattern:     "/ws",
		Port:        ":8443",
		TlsCertPath: "",
		TlsKeyPath:  "",
	}
}
