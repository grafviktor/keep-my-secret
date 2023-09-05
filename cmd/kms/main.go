package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v7"
	"golang.org/x/sync/errgroup"

	"github.com/grafviktor/keep-my-secret/internal/api/web"
	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/storage"
	"github.com/grafviktor/keep-my-secret/internal/version"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	// Display application version and build details
	version.Set(buildVersion, buildDate, buildCommit)
	version.PrintConsole()

	ec := config.EnvConfig{}
	if err := env.Parse(&ec); err != nil {
		log.Printf("%+v\n", err)
	}

	if ec.DevMode { // To enable dev mode, use 'DEV=true' env variable
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("=============================")
		fmt.Println("Running in Dev Mode !!!")
		fmt.Println("CORS is enabled for client connections")
		fmt.Println("Detailed logging is enabled")
		fmt.Println("=============================")
	}

	appConfig := config.New(ec)

	appContext, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		<-c
		cancel()
	}()

	dataStorage, err := storage.GetStorage(appContext, appConfig.StorageType, appConfig.DSN)
	if err != nil {
		log.Fatal(err)
	}

	router := web.NewHTTPRouter(appConfig, dataStorage)
	httpServer := http.Server{
		Addr:    appConfig.ServerAddr,
		Handler: router,
	}

	g, gCtx := errgroup.WithContext(appContext)
	g.Go(func() error {
		log.Printf("HTTP Server is listening on %s in secure mode\n", appConfig.ServerAddr)
		log.Printf("Client application is available at https://%s%s\n", appConfig.ServerAddr, appConfig.ClientAppURL)

		return httpServer.ListenAndServeTLS(appConfig.HTTPSCertPath, appConfig.HTTPSKeyPath)
	})

	g.Go(func() error {
		<-gCtx.Done()

		log.Println("Graceful shutdown")
		log.Println("Closing database connection")
		dataStorage.Close()

		log.Println("Shutting down web-server")
		return httpServer.Shutdown(appContext)
	})

	if err := g.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}

	log.Println("Exited")
}
