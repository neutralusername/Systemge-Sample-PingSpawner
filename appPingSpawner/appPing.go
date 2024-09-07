package appPingSpawner

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type AppPing struct {
	isStarted bool
	despawn   func()

	systemgeClient  *SystemgeClient.SystemgeClient
	dashboardClient *Dashboard.Client
}

func newAppPing(id string, despawn func()) *AppPing {
	app := &AppPing{
		despawn:   despawn,
		isStarted: true,
	}

	messageHandler := SystemgeConnection.NewConcurrentMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{
			"stop": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
				go app.close()
			},
		},
		SystemgeConnection.SyncMessageHandlers{
			"ping": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
				println("received ping request from", message.GetOrigin())
				return "", nil
			},
		},
		nil, nil,
	)
	app.systemgeClient = SystemgeClient.New(id,
		&Config.SystemgeClient{
			ClientConfigs: []*Config.TcpClient{
				{
					Address: "localhost:60001",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
			},
			ConnectionConfig: &Config.TcpSystemgeConnection{},
		},
		func(connection SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
		},
	)
	app.dashboardClient = Dashboard.NewClient(id,
		&Config.DashboardClient{
			ConnectionConfig: &Config.TcpSystemgeConnection{},
			ClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.systemgeClient.Start, app.close, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus,
		nil,
	)
	if err := app.dashboardClient.Start(); err != nil {
		panic(err)
	}
	if err := app.systemgeClient.Start(); err != nil {
		panic(err)
	}
	return app
}

func (app *AppPing) close() error {
	app.systemgeClient.Stop()
	app.dashboardClient.Stop()
	app.despawn()
	app.isStarted = false
	return nil
}
