package appWebsocketHTTP

import "sync/atomic"

type AppWebsocketHTTP struct {
	nextSpawnedNodePort *atomic.Uint32
}

func New() *AppWebsocketHTTP {
	port := &atomic.Uint32{}
	port.Store(60003)
	return &AppWebsocketHTTP{
		nextSpawnedNodePort: port,
	}
}
