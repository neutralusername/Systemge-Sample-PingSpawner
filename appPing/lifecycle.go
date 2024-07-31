package appPing

import (
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) OnStart(node *Node.Node) error {
	println(node.GetName() + " started")
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	println(node.GetName() + " ended")
	return nil
}
