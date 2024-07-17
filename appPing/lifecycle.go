package appPing

import (
	"Systemge/Node"
	"SystemgeSamplePingSpawner/topics"
)

func (app *App) OnStart(node *Node.Node) error {
	println(node.GetName() + " started")
	response, err := node.SyncMessage(topics.PING, node.GetName(), "ping")
	if err != nil {
		panic(err)
	}
	println(node.GetName() + " received \"" + response.GetPayload() + "\" from: " + response.GetOrigin())
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	println(node.GetName() + " ended")
	return nil
}
