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
	err := node.AsyncMessage(Spawner.SPAWN_AND_START_NODE_ASYNC, Helpers.JsonMarshal(&Config.NewNode{
		NodeConfig: &Config.Node{
			Name:                      "spawnedNode" + "-" + websocketClient.GetId(),
			RandomizerSeed:            Tools.GetSystemTime(),
			InfoLoggerPath:            "logs.log",
			WarningLoggerPath:         "logs.log",
			ErrorLoggerPath:           "logs.log",
			InternalInfoLoggerPath:    "logs.log",
			InternalWarningLoggerPath: "logs.log",
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
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:" + Helpers.IntToString(int(port)),
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
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
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Failed sending async message", err).Error())
		}
		websocketClient.Disconnect()
		return
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := node.AsyncMessage(Spawner.STOP_AND_DESPAWN_NODE_ASYNC, "spawnedNode"+"-"+websocketClient.GetId())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Failed sending async message", err).Error())
		}
		return
	}
}
