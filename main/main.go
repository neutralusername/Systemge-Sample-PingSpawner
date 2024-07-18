package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/TcpEndpoint"
	"Systemge/TcpServer"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/appWebsocketHTTP"
	"SystemgeSamplePingSpawner/config"
	"SystemgeSamplePingSpawner/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {

	Module.StartCommandLineInterface(Module.NewMultiModule(true,
		Node.New(Config.Node{
			Name:   "nodeResolver",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Resolver.New(Config.Resolver{
			Server:       TcpServer.New(60000, "MyCertificate.crt", "MyKey.key"),
			ConfigServer: TcpServer.New(60001, "MyCertificate.crt", "MyKey.key"),

			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerSpawner",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60002, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60002", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
			ConfigServer: TcpServer.New(60003, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "example.com", Utilities.GetFileContent("MyCertificate.crt")),

			SyncTopics:  []string{topics.END_NODE_SYNC, topics.START_NODE_SYNC},
			AsyncTopics: []string{topics.END_NODE_ASYNC, topics.START_NODE_ASYNC},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60004, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60004", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
			ConfigServer: TcpServer.New(60005, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "example.com", Utilities.GetFileContent("MyCertificate.crt")),

			SyncTopics: []string{topics.PINGPONG},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerPing",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60006, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60006", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
			ConfigServer: TcpServer.New(60007, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "example.com", Utilities.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeSpawner",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Spawner.New(Config.Spawner{
			SpawnedNodeLogger:      Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
			IsSpawnedNodeTopicSync: false,
			ResolverEndpoint:       TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			BrokerConfigEndpoint:   TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(60003), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
		}, Config.Systemge{
			HandleMessagesSequentially: false,

			BrokerSubscribeDelayMs:    1000,
			TopicResolutionLifetimeMs: 10000,
			SyncResponseTimeoutMs:     10000,
			TcpTimeoutMs:              5000,

			ResolverEndpoint: TcpEndpoint.New("127.0.0.1:60000", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
		},
			appPing.New)),
		Node.New(Config.Node{
			Name:   "nodeWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, appWebsocketHTTP.New()),
	))
}
