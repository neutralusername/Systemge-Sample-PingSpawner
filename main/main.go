package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/TcpEndpoint"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/appWebsocketHTTP"
	"SystemgeSamplePingSpawner/config"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	err := Resolver.New(Config.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerSpawner.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerPing.systemge")).Start()
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Node.New(Config.Node{
			Name:                      "nodePingSpawner",
			LoggerPath:                ERROR_LOG_FILE_PATH,
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			TopicResolutionLifetimeMs: 10000,
			BrokerReconnectDelayMs:    1000,
		}, Spawner.New(Config.Application{
			HandleMessagesSequentially: false,
		}, Config.Spawner{
			SpawnedNodeLoggerPath:  ERROR_LOG_FILE_PATH,
			IsSpawnedNodeTopicSync: true,
			ResolverEndpoint:       TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			BrokerConfigEndpoint:   TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.BROKER_CONFIG_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
		}, appPing.New)),
		Node.New(Config.Node{
			Name:                      "nodeWebsocketHTTP",
			LoggerPath:                ERROR_LOG_FILE_PATH,
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			TopicResolutionLifetimeMs: 10000,
			BrokerReconnectDelayMs:    1000,
		}, appWebsocketHTTP.New()),
	))
}
