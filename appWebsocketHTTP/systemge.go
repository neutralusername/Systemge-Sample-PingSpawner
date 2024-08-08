package appWebsocketHTTP

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		"spawnedNodeStarted": func(node *Node.Node, message *Message.Message) error {
			println(node.GetName() + " received spawnedNodeStarted")
			// these operation are done in order for this node to be able to message the spawned node
			tcpEndpointConfig := Config.UnmarshalTcpEndpoint(message.GetPayload())
			node.ConnectToNode(tcpEndpointConfig, false)

			// ping-pong to check if connection is working
			responseChannel, err := node.SyncMessage("ping", "")
			println(node.GetName() + " sent ping-sync")
			if err != nil {
				panic(err)
			}
			_, err = responseChannel.ReceiveResponse()
			if err != nil {
				panic(err)
			}
			println(node.GetName() + " received pong-sync")
			return nil
		},
		"SpawnedNodeStopped": func(node *Node.Node, message *Message.Message) error {
			println(node.GetName() + " received SpawnedNodeStopped")
			// these operations are done in order to stop the reconnection-attempts to now stopped node
			tcpEndpointConfig := Config.UnmarshalTcpEndpoint(message.GetPayload())
			node.DisconnectFromNode(tcpEndpointConfig.Address)
			return nil
		},
		"ping": func(node *Node.Node, message *Message.Message) error {
			println(node.GetName() + " received ping-async")
			err := node.AsyncMessage("pong", "")
			if err != nil {
				panic(err)
			}
			println(node.GetName() + " sent pong-async")
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}
