package appWebsocketHttp

import (
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Status"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
	"github.com/neutralusername/Systemge/SystemgeServer"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	status      int
	statusMutex sync.Mutex

	systemgeServer  *SystemgeServer.SystemgeServer
	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}

	messageHandler := SystemgeMessageHandler.New(
		SystemgeMessageHandler.AsyncMessageHandlers{
			"ping": func(message *Message.Message) {
				println("received ping-async")
				err := app.systemgeServer.AsyncMessage("pong", "", message.GetOrigin())
				if err != nil {
					panic(err)
				}
				println("sent pong-async")
			},
		},
		SystemgeMessageHandler.SyncMessageHandlers{},
	)
	app.systemgeServer = SystemgeServer.New(
		&Config.SystemgeServer{
			Name:              "systemgeServer",
			InfoLoggerPath:    "logs.log",
			WarningLoggerPath: "logs.log",
			ErrorLoggerPath:   "logs.log",
			ListenerConfig: &Config.SystemgeListener{
				TcpListenerConfig: &Config.TcpListener{
					TlsCertPath: "MyCertificate.crt",
					TlsKeyPath:  "MyKey.key",
					Port:        60001,
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		func(connection *SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			switch connection.GetName() {
			case "appSpawner":
				return nil
			default:
				println(connection.GetName() + " connected")
				responseChannel, err := app.systemgeServer.SyncRequest("ping", "", connection.GetName())
				if err != nil {
					panic(err)
				}
				println("sent ping request to " + connection.GetName())
				response := <-responseChannel
				if response == nil {
					panic(Error.New("response is nil", nil))
				}
				println("received ping response from " + connection.GetName())
				return nil
			}
		},
		func(connection *SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
			println(connection.GetName() + " disconnected")
		},
	)
	app.websocketServer = WebsocketServer.New(
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpListenerConfig: &Config.TcpListener{
				Port: 8443,
			},
		},
		WebsocketServer.MessageHandlers{},
		app.OnConnectHandler, app.OnDisconnectHandler,
	)
	app.httpServer = HTTPServer.New(
		&Config.HTTPServer{
			TcpListenerConfig: &Config.TcpListener{
				Port: 8080,
			},
		},
		HTTPServer.Handlers{
			"/": HTTPServer.SendDirectory("../frontend"),
		},
	)
	Dashboard.NewClient(
		&Config.DashboardClient{
			Name:             "appWebsocketHttp",
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.start, app.stop, app.systemgeServer.GetMetrics, app.getStatus,
		nil,
	)
	return app
}

func (app *AppWebsocketHTTP) getStatus() int {
	return app.status
}

func (app *AppWebsocketHTTP) start() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STOPPED {
		return Error.New("App already started", nil)
	}
	if err := app.systemgeServer.Start(); err != nil {
		return Error.New("Failed to start systemgeServer", err)
	}
	if err := app.websocketServer.Start(); err != nil {
		app.systemgeServer.Stop()
		return Error.New("Failed to start websocketServer", err)
	}
	if err := app.httpServer.Start(); err != nil {
		app.systemgeServer.Stop()
		app.websocketServer.Stop()
		return Error.New("Failed to start httpServer", err)
	}
	app.status = Status.STARTED
	return nil
}

func (app *AppWebsocketHTTP) stop() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STARTED {
		return Error.New("App not started", nil)
	}
	app.httpServer.Stop()
	app.websocketServer.Stop()
	app.systemgeServer.Stop()
	app.status = Status.STOPPED
	return nil
}

func (app *AppWebsocketHTTP) WebsocketPropagate(message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	responseChannel, err := app.systemgeServer.SyncRequest("spawn", websocketClient.GetId(), "appSpawner")
	if err != nil {
		panic(err)
	}
	response := <-responseChannel
	if response == nil {
		panic(Error.New("response is nil", nil))
	}
	if response.GetTopic() == Message.TOPIC_FAILURE {
		panic(Error.New("response is failure", nil))
	}
	println("test")
	return nil
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(websocketClient *WebsocketServer.WebsocketClient) {
	app.systemgeServer.AsyncMessage("stop", "", websocketClient.GetId())
}
