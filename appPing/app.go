package appPing

import (
	"github.com/neutralusername/Systemge/Node"
)

type App struct {
	isStarted bool
}

func New() Node.Application {
	return &App{}
}
