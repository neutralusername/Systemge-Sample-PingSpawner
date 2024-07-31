package appPing

import (
	"github.com/neutralusername/Systemge/Node"
)

type App struct {
}

func New() Node.Application {
	return &App{}
}
