package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"groupie-tracker/internal/controller/router"
	"groupie-tracker/internal/controller/server"
	"groupie-tracker/pkg/config"
)

// Run ...
func Run(cfg *config.Config) error {
	// init router +
	router := router.New()
	if err := router.InitRouter(); err != nil {
		slog.Error("error occured while the ServerMux router initialization")
		return err
	}
	slog.Debug("the ServeMux router successful created and initialized")

	//http server +
	srv := server.New(cfg, router)

	slog.Debug("http server successfully created")

	//running the http server + -> for{ infinity logic}
	fmt.Printf("\n\nServer running on -> http://%s:%s\n\n", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	go func() {
		if err := srv.Run(); err != nil {
			slog.Error(err.Error())
		}
	}()

	// getting API -> webapi <-controller ---

	// gracefull shutdown -> stopping the http server +
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// logging INFO
	sig := <-sigChan
	slog.Info("Some crazy developer termination our app, with: ", "syscall", sig.String())

	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error(err.Error())

		return err
	}

	return nil
}

