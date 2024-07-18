package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/Tcp"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/appWebsocketHTTP"
	"SystemgeSamplePingSpawner/config"
	"SystemgeSamplePingSpawner/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {

	Node.StartCommandLineInterface(true,
		Node.New(Config.Node{
			Name: "nodeResolver",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Resolver.New(Config.Resolver{
			Server:       Tcp.NewServer(60000, "MyCertificate.crt", "MyKey.key"),
			ConfigServer: Tcp.NewServer(60001, "MyCertificate.crt", "MyKey.key"),

			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerSpawner",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60002, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60002", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60003, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncTopics:  []string{topics.END_NODE_SYNC, topics.START_NODE_SYNC},
			AsyncTopics: []string{topics.END_NODE_ASYNC, topics.START_NODE_ASYNC},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60004, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60004", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60005, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncTopics: []string{topics.PINGPONG},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerPing",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60006, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60006", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60007, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeSpawner",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Spawner.New(Config.Spawner{
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
			IsSpawnedNodeTopicSync: false,
			ResolverEndpoint:       Tcp.NewEndpoint(config.SERVER_IP+":"+Helpers.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Helpers.GetFileContent(config.CERT_PATH)),
			BrokerConfigEndpoint:   Tcp.NewEndpoint(config.SERVER_IP+":"+Helpers.IntToString(60003), config.SERVER_NAME_INDICATION, Helpers.GetFileContent(config.CERT_PATH)),
		}, Config.Systemge{
			HandleMessagesSequentially: false,

			BrokerSubscribeDelayMs:    1000,
			TopicResolutionLifetimeMs: 10000,
			SyncResponseTimeoutMs:     10000,
			TcpTimeoutMs:              5000,

			ResolverEndpoint: Tcp.NewEndpoint("127.0.0.1:60000", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
		},
			appPing.New)),
		Node.New(Config.Node{
			Name: "nodeWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, appWebsocketHTTP.New()),
	)
}
