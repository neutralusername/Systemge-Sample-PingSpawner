package main

import (
	"Systemge/Module"
	"SystemgeSamplePingSpawner/appSpawner"
	"SystemgeSamplePingSpawner/appWebsocketHTTP"
)

const TOPICRESOLUTIONSERVER_ADDRESS = "127.0.0.1:60000"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	err := Module.NewResolverFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerSpawner.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerPing.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	clientSpawner := Module.NewClient("clientSpawner", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appSpawner.New, nil)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Module.NewCompositeClientWebsocketHTTP("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", HTTP_PORT, "", "", appWebsocketHTTP.New, nil),
		clientSpawner,
	), clientSpawner.GetApplication().GetCustomCommandHandlers())
}
