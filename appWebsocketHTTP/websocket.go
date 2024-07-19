package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
	"net/http"

	"github.com/gorilla/websocket"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	_, err := node.SyncMessage(topics.START_NODE_SYNC, websocketClient.GetId(), websocketClient.GetId())
	if err != nil {
		panic(Error.New("Error sending sync message", err))
	}
	err = node.AsyncMessage(websocketClient.GetId(), websocketClient.GetId(), "ping")
	if err != nil {
		panic(Error.New("Error sending async message", err))
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	_, err := node.SyncMessage(topics.END_NODE_SYNC, node.GetName(), websocketClient.GetId())
	if err != nil {
		//windows seems to have issues with the sync token generation.. sometimes it will generate two similar tokens in sequence. i assume the system time is not accurate enough for very fast token generation
		panic(Error.New("Error sending sync message", err))
	}
}

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() *Config.Websocket {
	return &Config.Websocket{
		Pattern: "/ws",
		Server: &Config.TcpServer{
			Port: 8443,
		},
		HandleClientMessagesSequentially: false,
		ClientMessageCooldownMs:          0,
		ClientWatchdogTimeoutMs:          20000,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
