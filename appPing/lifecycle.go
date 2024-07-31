package appPing

import (
	"time"

	"github.com/neutralusername/Systemge/Node"
)

func (app *App) OnStart(node *Node.Node) error {
	app.isStarted = true
	println(node.GetName() + " started")
	go func() {
		for app.isStarted {
			node.AsyncMessage("ping", node.GetName())
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	app.isStarted = false
	println(node.GetName() + " ended")
	return nil
}
