package appWebsocketHTTP

import (
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		"spawnedNodeStarted": func(node *Node.Node, message *Message.Message) error {
			// these operation are done in order for this node to be able to message the spawned node
			tcpEndpointConfig := Config.UnmarshalTcpEndpoint(message.GetPayload())
			node.StartOutgoingConnectionLoop(tcpEndpointConfig)

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
			// these operations are done in order to stop the reconnection-attempts to now stopped node
			// there is currently the issue, that this message is received while the connection is still active, which means cancelling the ongoing connection loop will do nothing
			// temporary fix is either a delay or limit the reconnection attempts
			// but i will try to fix this today
			time.Sleep(1 * time.Second)
			tcpEndpointConfig := Config.UnmarshalTcpEndpoint(message.GetPayload())
			node.CancelOutgoingConnectionLoop(tcpEndpointConfig.Address)
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}
