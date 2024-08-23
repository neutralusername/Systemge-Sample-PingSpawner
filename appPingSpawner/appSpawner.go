package appPingSpawner

import (
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
)

type AppSpawner struct {
	spawnedApps map[string]*AppPing
	mutex       *sync.Mutex

	systemgeClient *SystemgeClient.SystemgeClient
}

func New() *AppSpawner {
	app := &AppSpawner{
		spawnedApps: make(map[string]*AppPing),
		mutex:       &sync.Mutex{},
	}

	app.systemgeClient = SystemgeClient.New(
		&Config.SystemgeClient{
			Name:              "appSpawner",
			InfoLoggerPath:    "logs.log",
			WarningLoggerPath: "logs.log",
			ErrorLoggerPath:   "logs.log",
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
			SystemgeMessageHandler.AsyncMessageHandlers{},
			SystemgeMessageHandler.SyncMessageHandlers{
				"spawn": func(message *Message.Message) (string, error) {
					app.mutex.Lock()
					defer app.mutex.Unlock()

					if _, ok := app.spawnedApps[message.GetPayload()]; ok {
						return "", Error.New("app \""+message.GetPayload()+"\" already spawned", nil)
					}
					app.spawnedApps[message.GetPayload()] = newAppPing(message.GetPayload(), func() {
						app.mutex.Lock()
						defer app.mutex.Unlock()

						delete(app.spawnedApps, message.GetPayload())
					})
					return "", nil
				},
			},
		),
	)
	Dashboard.NewClient(&Config.DashboardClient{
		Name:             "appSpawner",
		ConnectionConfig: &Config.SystemgeConnection{},
		EndpointConfig: &Config.TcpEndpoint{
			Address: "localhost:60000",
			TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			Domain:  "example.com",
		},
	}, app.systemgeClient.Start, app.systemgeClient.Stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus, nil)

	return app
}
