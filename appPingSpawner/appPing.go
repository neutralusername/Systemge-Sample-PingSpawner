package appPingSpawner

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardClientCustomService"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type AppPing struct {
	isStarted bool
	despawn   func()

	systemgeClient  *SystemgeClient.SystemgeClient
	dashboardClient *DashboardClientCustomService.Client
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
			TcpClientConfigs: []*Config.TcpClient{
				{
					Address: "localhost:60001",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
			},
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
		},
		func(connection SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
		},
	)
	app.dashboardClient = DashboardClientCustomService.New(id,
		&Config.DashboardClient{
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			TcpClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.systemgeClient,
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
