package appSpawner

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSamplePingSpawner/appPing"
	"SystemgeSamplePingSpawner/topics"
	"sync"
)

type App struct {
	client *Client.Client

	// the games topic resolver will be on port 6002x
	chessRoomIds map[string]bool           // x -> inUse
	chessClients map[string]*Client.Client // x -> client
	mutex        sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client: client,

		chessRoomIds: map[string]bool{
			"0": false,
			"1": false,
			"2": false,
			"3": false,
			"4": false,
			"5": false,
			"6": false,
			"7": false,
			"8": false,
			"9": false,
		},
		chessClients: map[string]*Client.Client{},
	}
	return app, nil
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.NEW: app.NewGame,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *App) NewGame(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	availableId := ""
	for port, inUse := range app.chessRoomIds {
		if !inUse {
			availableId = port
			break
		}
	}
	if availableId == "" {
		return "", Utilities.NewError("No available room numbers", nil)
	}
	moveTopic := "move_" + availableId
	pingClient := Client.New("clientPing"+availableId, app.client.GetTopicResolutionServerAddress(), app.client.GetLogger(), nil)
	pingApp, err := appPing.New(pingClient, []string{moveTopic})
	if err != nil {
		return "", Utilities.NewError("Error creating ping app "+availableId, err)
	}
	pingClient.SetApplication(pingApp)
	brokerNetConn, err := Utilities.TlsDial("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return "", Utilities.NewError("Error dialing ping broker", err)
	}
	_, err = Utilities.TcpExchange(brokerNetConn, Message.NewAsync("addAsyncTopic", app.client.GetName(), moveTopic), 5000)
	if err != nil {
		return "", Utilities.NewError("Error exchanging messages with chess broker", err)
	}
	resolverNetConn, err := Utilities.TlsDial("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return "", Utilities.NewError("Error dialing topic resolution server", err)
	}
	_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("registerTopics", app.client.GetName(), "brokerPing "+moveTopic), 5000)
	if err != nil {
		return "", Utilities.NewError("Error exchanging messages with topic resolution server", err)
	}
	err = pingClient.Start()
	if err != nil {
		return "", Utilities.NewError("Error starting chess client", err)
	}

	app.chessRoomIds[availableId] = true
	app.chessClients[availableId] = pingClient
	return "", nil
}
