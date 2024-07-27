package appPing

import (
	"github.com/neutralusername/Systemge/Node"
)

type App struct {
	id string
}

func New(id string) Node.Application {
	app := &App{
		id: id,
	}
	return app
}
