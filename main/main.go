package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"SystemgeSamplePingSpawner/appSpawner"
	"SystemgeSamplePingSpawner/appWebsocketHTTP"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	err := Resolver.New(Module.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Module.ParseBrokerConfigFromFile("brokerSpawner.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Module.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Module.ParseBrokerConfigFromFile("brokerPing.systemge")).Start()
	if err != nil {
		panic(err)
	}
	applicationWebsocketHTTP := appWebsocketHTTP.New()
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Node.New(Config.Node{
			Name:       "nodeWebsocketHTTP",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, applicationWebsocketHTTP, applicationWebsocketHTTP, applicationWebsocketHTTP),
		Node.New(Config.Node{
			Name:       "nodeSpawner",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, appSpawner.New(), nil, nil),
	))
}
