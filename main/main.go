package main

import (
	"SystemgeSamplePingSpawner/appPingSpawner"
	"SystemgeSamplePingSpawner/appWebsocketHttp"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
)

const LOGGER_PATH = "logs.log"

func main() {
	Dashboard.NewServer(&Config.DashboardServer{
		InfoLoggerPath:    LOGGER_PATH,
		WarningLoggerPath: LOGGER_PATH,
		ErrorLoggerPath:   LOGGER_PATH,
		HTTPServerConfig: &Config.HTTPServer{
			TcpListenerConfig: &Config.TcpListener{
				Port: 8081,
			},
		},
		WebsocketServerConfig: &Config.WebsocketServer{
			Pattern:                 "/ws",
			ClientWatchdogTimeoutMs: 1000 * 60,
			TcpListenerConfig: &Config.TcpListener{
				Port: 8444,
			},
		},
		SystemgeServerConfig: &Config.SystemgeServer{
			Name:              "dashboardServer",
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
			ListenerConfig: &Config.SystemgeListener{
				TcpListenerConfig: &Config.TcpListener{
					TlsCertPath: "MyCertificate.crt",
					TlsKeyPath:  "MyKey.key",
					Port:        60000,
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		HeapUpdateIntervalMs:      1000,
		GoroutineUpdateIntervalMs: 1000,
		StatusUpdateIntervalMs:    1000,
		MetricsUpdateIntervalMs:   1000,
	})
	appWebsocketHttp.New()
	appPingSpawner.New()
	<-make(chan time.Time)
}
