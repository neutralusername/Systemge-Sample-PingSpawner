package appPingSpawner

import (
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
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
		nil, nil,
		SystemgeMessageHandler.New(
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
		),
	)
	app.dashboardClient = Dashboard.NewClient(&Config.DashboardClient{
		Name:             id,
		ConnectionConfig: &Config.SystemgeConnection{},
		EndpointConfig: &Config.TcpEndpoint{
			Address: "localhost:60000",
			TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			Domain:  "example.com",
		},
	}, app.systemgeClient.Start, app.stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus, nil)

	err := app.systemgeClient.Start()
	if err != nil {
		panic(err)
	}
	return app
}

func (app *AppPing) stop() error {
	println("t1")
	app.systemgeClient.Stop()
	println("t2")
	app.dashboardClient.Close()
	println("t3")
	app.despawn()
	println("t4")
	app.isStarted = false
	return nil
}
