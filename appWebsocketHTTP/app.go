package appWebsocketHTTP

import (
	"sync"
	"sync/atomic"

	"github.com/neutralusername/Systemge/Config"
)

type AppWebsocketHTTP struct {
	nextSpawnedNodePort *atomic.Uint32
	activePorts         map[string]*Config.TcpEndpoint // websocketClient.GetId() -> port (for onDisconnectHandler, to stop re-connection attempts to despawned nodes)
	mutex               sync.Mutex
}

func New() *AppWebsocketHTTP {
	port := &atomic.Uint32{}
	port.Store(60003)
	return &AppWebsocketHTTP{
		nextSpawnedNodePort: port,
		activePorts:         make(map[string]*Config.TcpEndpoint),
	}
}
