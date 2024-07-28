package main

import (
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/appWebsocketHTTP"
	"SystemgeSamplePingSpawner/topics"

	"github.com/neutralusername/Systemge/Broker"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Resolver"
	"github.com/neutralusername/Systemge/Spawner"
	"github.com/neutralusername/Systemge/Tools"
)

const LOGGER_PATH = "logs.log"

func main() {
	loggerQueue := Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
	}, Dashboard.New(&Config.Dashboard{
		Server: &Config.TcpServer{
			Port: 8081,
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		NodeBrokerCounterIntervalMs:    1000,
		NodeResolverCounterIntervalMs:  1000,
		HeapUpdateIntervalMs:           1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		AutoStart:                      true,
		EnableDashboardCounters:        true,
	},
		Node.New(&Config.Node{
			Name:                  "nodeResolver",
			RandomizerSeed:        Tools.GetSystemTime(),
			InfoLogger:            Tools.NewLogger("[Info \"nodeResolver\"]", loggerQueue),
			InternalInfoLogger:    Tools.NewLogger("[InternalInfo \"nodeResolver\"]", loggerQueue),
			WarningLogger:         Tools.NewLogger("[Warning \"nodeResolver\"] ", loggerQueue),
			InternalWarningLogger: Tools.NewLogger("[InternalWarning \"nodeResolver\"] ", loggerQueue),
			ErrorLogger:           Tools.NewLogger("[Error \"nodeResolver\"] ", loggerQueue),
		}, Resolver.New(&Config.Resolver{
			Server: &Config.TcpServer{
				Port:        60000,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ConfigServer: &Config.TcpServer{
				Port:        60001,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},

			TcpTimeoutMs: 5000,
		})),
		Node.New(&Config.Node{
			Name:                  "nodeBrokerSpawner",
			RandomizerSeed:        Tools.GetSystemTime(),
			InfoLogger:            Tools.NewLogger("[Info \"nodeBrokerSpawner\"]", loggerQueue),
			InternalInfoLogger:    Tools.NewLogger("[InternalInfo \"nodeBrokerSpawner\"]", loggerQueue),
			WarningLogger:         Tools.NewLogger("[Warning \"nodeBrokerSpawner\"] ", loggerQueue),
			InternalWarningLogger: Tools.NewLogger("[InternalWarning \"nodeBrokerSpawner\"] ", loggerQueue),
			ErrorLogger:           Tools.NewLogger("[Error \"nodeBrokerSpawner\"] ", loggerQueue),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60002,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60002",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{
				Port:        60003,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncTopics:  []string{Spawner.STOP_NODE_SYNC, Spawner.START_NODE_SYNC, Spawner.SPAWN_NODE_SYNC, Spawner.DESPAWN_NODE_SYNC},
			AsyncTopics: []string{Spawner.STOP_NODE_ASYNC, Spawner.START_NODE_ASYNC, Spawner.SPAWN_NODE_ASYNC, Spawner.DESPAWN_NODE_ASYNC},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:                  "nodeBrokerWebsocketHTTP",
			RandomizerSeed:        Tools.GetSystemTime(),
			InfoLogger:            Tools.NewLogger("[Info \"nodeBrokerWebsocketHTTP\"]", loggerQueue),
			InternalInfoLogger:    Tools.NewLogger("[InternalInfo \"nodeBrokerWebsocketHTTP\"]", loggerQueue),
			WarningLogger:         Tools.NewLogger("[Warning \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
			InternalWarningLogger: Tools.NewLogger("[InternalWarning \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:           Tools.NewLogger("[Error \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60004,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60004",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{
				Port:        60005,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncTopics: []string{topics.PINGPONG},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:                  "nodeBrokerPing",
			RandomizerSeed:        Tools.GetSystemTime(),
			InfoLogger:            Tools.NewLogger("[Info \"nodeBrokerPing\"]", loggerQueue),
			InternalInfoLogger:    Tools.NewLogger("[InternalInfo \"nodeBrokerPing\"]", loggerQueue),
			WarningLogger:         Tools.NewLogger("[Warning \"nodeBrokerPing\"] ", loggerQueue),
			InternalWarningLogger: Tools.NewLogger("[InternalWarning \"nodeBrokerPing\"] ", loggerQueue),
			ErrorLogger:           Tools.NewLogger("[Error \"nodeBrokerPing\"] ", loggerQueue),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60006,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60006",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{Port: 60007, TlsCertPath: "MyCertificate.crt", TlsKeyPath: "MyKey.key"},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:                  "nodeSpawner",
			RandomizerSeed:        Tools.GetSystemTime(),
			InfoLogger:            Tools.NewLogger("[Info \"nodeSpawner\"]", loggerQueue),
			InternalInfoLogger:    Tools.NewLogger("[InternalInfo \"nodeSpawner\"]", loggerQueue),
			WarningLogger:         Tools.NewLogger("[Warning \"nodeSpawner\"] ", loggerQueue),
			InternalWarningLogger: Tools.NewLogger("[InternalWarning \"nodeSpawner\"] ", loggerQueue),
			ErrorLogger:           Tools.NewLogger("[Error \"nodeSpawner\"] ", loggerQueue),
		}, Spawner.New(&Config.Spawner{
			LoggerQueue:                 loggerQueue,
			IsSpawnedNodeTopicSync:      false,
			PropagateSpawnedNodeChanges: true,
			ResolverEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60000",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			BrokerConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60003",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
		}, &Config.Systemge{
			HandleMessagesSequentially: false,

			BrokerSubscribeDelayMs:    1000,
			TopicResolutionLifetimeMs: 10000,
			SyncResponseTimeoutMs:     10000,
			TcpTimeoutMs:              5000,

			ResolverEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60000",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
		},
			appPing.New)),
		Node.New(&Config.Node{
			Name:                  "nodeWebsocketHTTP",
			RandomizerSeed:        Tools.GetSystemTime(),
			InfoLogger:            Tools.NewLogger("[Info \"nodeWebsocketHTTP\"]", loggerQueue),
			InternalInfoLogger:    Tools.NewLogger("[InternalInfo \"nodeWebsocketHTTP\"]", loggerQueue),
			WarningLogger:         Tools.NewLogger("[Warning \"nodeWebsocketHTTP\"] ", loggerQueue),
			InternalWarningLogger: Tools.NewLogger("[InternalWarning \"nodeWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:           Tools.NewLogger("[Error \"nodeWebsocketHTTP\"] ", loggerQueue),
		}, appWebsocketHTTP.New()),
	)).StartBlocking()
}
