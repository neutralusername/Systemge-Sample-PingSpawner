package appPing

import (
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) OnStart(node *Node.Node) error {
	err := node.AsyncMessage("spawnedNodeStarted", Helpers.JsonMarshal(node.GetSystemgeEndpointConfig()))
	if err != nil {
		return err
	}
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	err := node.AsyncMessage("SpawnedNodeStopped", Helpers.JsonMarshal(node.GetSystemgeEndpointConfig()))
	if err != nil {
		return err
	}
	return nil
}
