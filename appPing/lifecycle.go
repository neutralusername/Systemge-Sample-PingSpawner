package appPing

import (
	"time"

	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) OnStart(node *Node.Node) error {
	err := node.AsyncMessage("spawnedNodeStarted", Helpers.JsonMarshal(node.GetSystemgeEndpointConfig()))
	if err != nil {
		return err
	}
	println(node.GetName() + " sent spawnedNodeStarted")
	app.isStarted = true
	go func() {
		for app.isStarted {
			err := node.AsyncMessage("ping", node.GetName())
			if err != nil {
				panic(err)
			}
			println(node.GetName() + " sent ping-async")
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	app.isStarted = false
	err := node.AsyncMessage("SpawnedNodeStopped", Helpers.JsonMarshal(node.GetSystemgeEndpointConfig()))
	if err != nil {
		return err
	}
	println(node.GetName() + " sent SpawnedNodeStopped")
	return nil
}
