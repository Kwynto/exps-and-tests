package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api/internal/config"
	"rest-api/internal/user"
	"rest-api/pkg/simlog"
	"time"
)

func start(router *http.ServeMux, cfg *config.Config) {
	simlog.Info("Start application.")

	var listener net.Listener
	var errList error

	if cfg.Listen.Type == "sock" {
		appDir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
		if errDir != nil {
			simlog.Fatal(errDir)
		}

		simlog.Debug("create socket")
		socketPath := path.Join(appDir, "app.sock")
		simlog.Debug(fmt.Sprintf("socket path: %s", socketPath))

		simlog.Debug("listen unix socket")
		listener, errList = net.Listen("unix", socketPath)
	} else {
		hostPort := fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
		simlog.Debugf("listen tcp %s", hostPort)
		listener, errList = net.Listen("tcp", hostPort)
	}

	if errList != nil {
		simlog.Critical("Not create listener: %v", errList)
		simlog.Fatal(errList)
	}

	simlog.Debug("create structure of server")
	server := http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     simlog.GetStdLogger(),
	}

	simlog.Info("Start server.")
	servErr := server.Serve(listener)
	if servErr != nil {
		simlog.Printf("Error: %v \n", servErr)
		return
	}
	simlog.Info("Stop server.")
}

func main() {
	simlog.Info("Reading configuration")
	cfg := config.GetConfig()

	simlog.Info("Create router.")
	router := http.NewServeMux()
	user.NewHandler().Register(router)

	start(router, cfg)
	simlog.Info("Program finished.")
}
