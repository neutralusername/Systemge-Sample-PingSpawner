package appPingSpawner

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
)

type AppPing struct {
	isStarted bool
	despawn   func()

	systemgeClient  *SystemgeClient.SystemgeClient
	dashboardClient *Dashboard.DashboardClient
}

func newAppPing(id string, despawn func()) *AppPing {
	app := &AppPing{
		despawn:   despawn,
		isStarted: true,
	}

	messageHandler := SystemgeMessageHandler.NewConcurrentMessageHandler(
		SystemgeMessageHandler.AsyncMessageHandlers{
			"stop": func(message *Message.Message) {
				go app.stop()
			},
		},
		SystemgeMessageHandler.SyncMessageHandlers{
			"ping": func(message *Message.Message) (string, error) {
				println("received ping request from", message.GetOrigin())
				return "", nil
			},
		},
		nil, nil,
	)
	app.systemgeClient = SystemgeClient.New(
		&Config.SystemgeClient{
			Name: id,
			EndpointConfigs: []*Config.TcpEndpoint{
				{
					Address: "localhost:60001",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		func(connection *SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			return nil
		},
		func(connection *SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
		},
	)
	app.dashboardClient = Dashboard.NewClient(
		&Config.DashboardClient{
			Name:             id,
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.systemgeClient.Start, app.stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus,
		nil,
	)

	if err := app.systemgeClient.Start(); err != nil {
		panic(err)
	}
	return app
}

func (app *AppPing) stop() error {
	app.systemgeClient.Stop()
	app.dashboardClient.Close()
	app.despawn()
	app.isStarted = false
	return nil
}
