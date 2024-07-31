package appWebsocketHTTP

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Spawner"
	"github.com/neutralusername/Systemge/Tools"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	port := app.nextSpawnedNodePort.Add(1)
	_, err := node.SyncMessage(Spawner.SPAWN_NODE_SYNC, Helpers.JsonMarshal(&Config.NewNode{
		NodeConfig: &Config.Node{
			Name:                      "spawnedNode" + "-" + websocketClient.GetId(),
			RandomizerSeed:            Tools.GetSystemTime(),
			InfoLoggerPath:            "logs.log",
			WarningLoggerPath:         "logs.log",
			ErrorLoggerPath:           "logs.log",
			InternalInfoLoggerPath:    "logs.log",
			InternalWarningLoggerPath: "logs.log",
			/* 	InternalInfoLoggerPath:    "logs.log",
			InternalWarningLoggerPath: "logs.log", */
		},
		SystemgeConfig: &Config.Systemge{
			HandleMessagesSequentially: false,

			SyncRequestTimeoutMs:            10000,
			TcpTimeoutMs:                    5000,
			MaxConnectionAttempts:           0,
			ConnectionAttemptDelayMs:        1000,
			StopAfterOutgoingConnectionLoss: true,
			ServerConfig: &Config.TcpServer{
				Port:        uint16(port),
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			EndpointConfigs: []*Config.TcpEndpoint{
				{
					Address: "localhost:60001",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
				{
					Address: "localhost:60002",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
			},
			IncomingMessageByteLimit: 0,
			MaxPayloadSize:           0,
			MaxTopicSize:             0,
			MaxSyncTokenSize:         0,
			SyncResponseLimit:        1,
		},
	}))
	if err != nil {
		panic(Error.New("Failed sending sync message", err))
	}
	_, err = node.SyncMessage(Spawner.START_NODE_SYNC, "spawnedNode"+"-"+websocketClient.GetId())
	if err != nil {
		panic(Error.New("Failed sending sync message", err))
	}
	tcpEndpointConfig := &Config.TcpEndpoint{
		Address: "localhost:" + Helpers.IntToString(int(port)),
		TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
		Domain:  "example.com",
	}
	app.mutex.Lock()
	app.activePorts["spawnedNode"+"-"+websocketClient.GetId()] = tcpEndpointConfig
	app.mutex.Unlock()
	node.OutgoingConnectionLoop(tcpEndpointConfig)
	responseChannel, err := node.SyncMessage("ping", "")
	println(node.GetName() + " sent ping-sync")
	if err != nil {
		panic(Error.New("Failed sending sync message", err))
	}
	_, err = responseChannel.ReceiveResponse()
	if err != nil {
		panic(Error.New("Failed receiving response", err))
	}
	println(node.GetName() + " received pong-sync")
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	responseChannel, err := node.SyncMessage(Spawner.STOP_NODE_SYNC, "spawnedNode"+"-"+websocketClient.GetId())
	if err != nil {
		panic(Error.New("Failed sending sync message", err))
	}
	_, err = responseChannel.ReceiveResponse()
	if err != nil {
		panic(Error.New("Failed receiving response", err))
	}
	err = node.AsyncMessage(Spawner.DESPAWN_NODE_ASYNC, "spawnedNode"+"-"+websocketClient.GetId())
	if err != nil {
		panic(Error.New("Failed sending async message", err))
	}
	app.mutex.Lock()
	tcpEndpointConfig := app.activePorts["spawnedNode"+"-"+websocketClient.GetId()]
	delete(app.activePorts, "spawnedNode"+"-"+websocketClient.GetId())
	app.mutex.Unlock()
	node.CancelOutgoingConnectionLoop(tcpEndpointConfig.Address)
}
