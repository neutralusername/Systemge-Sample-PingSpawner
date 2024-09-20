package appPingSpawner

import (
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardClientCustomService"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
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

	messageHandler := SystemgeConnection.NewConcurrentMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{},
		SystemgeConnection.SyncMessageHandlers{
			"spawn": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
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
		nil, nil,
	)
	app.systemgeClient = SystemgeClient.New("appSpawner",
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
			connection.StartMessageHandlingLoop_Sequentially(messageHandler)
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			connection.StopMessageHandlingLoop()
		},
	)
	DashboardClientCustomService.New("appSpawner",
		&Config.DashboardClient{
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			TcpClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		}, app.systemgeClient,
		nil,
	).Start()

	return app
}
