package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{}
}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	_, err := client.SyncMessage(topics.NEW, websocketClient.GetId(), websocketClient.GetId())
	if err != nil {
		panic(Error.New("Error sending sync message", err))
	}
	err = client.AsyncMessage(websocketClient.GetId(), websocketClient.GetId(), "ping")
	if err != nil {
		panic(Error.New("Error sending async message", err))
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	_, err := client.SyncMessage(topics.END, client.GetName(), websocketClient.GetId())
	if err != nil {
		//windows seems to have issues with the sync token generation.. sometimes it will generate two similar tokens in sequence. i assume the system time is not accurate enough for very fast token generation
		panic(Error.New("Error sending sync message", err))
	}
	client.RemoveTopicResolution(websocketClient.GetId())
}
