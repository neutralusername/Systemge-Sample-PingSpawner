package appPingSpawner

import (
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
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
			Name: "appSpawner",
			EndpointConfigs: []*Config.TcpEndpoint{
				{
					Address: "localhost:60001",
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
					app.spawnedApps[message.GetPayload()] = newAppPing(message.GetPayload())
					return "", nil
				},
				"despawn": func(message *Message.Message) (string, error) {
					app.mutex.Lock()
					defer app.mutex.Unlock()

					pingApp := app.spawnedApps[message.GetPayload()]
					if pingApp == nil {
						return "", Error.New("app \""+message.GetPayload()+"\" not spawned", nil)
					}
					pingApp.stop()
					delete(app.spawnedApps, message.GetPayload())
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
		},
	}, app.start, app.stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus, nil)

	return app
}

func (app *AppSpawner) start() error {
	err := app.systemgeClient.Start()
	if err != nil {
		return err
	}
	return nil
}

func (app *AppSpawner) stop() error {
	if err := app.systemgeClient.Stop(); err != nil {
		return err
	}
	return nil
}
