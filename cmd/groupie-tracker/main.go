package main

import (
	"log/slog"

	"groupie-tracker/pkg/config"
	"groupie-tracker/pkg/logger"
)

func main() {
	// init config ....
	cfg, err := config.InitConfig()
	if err != nil {
		// loggining -> Error or Fatal level
		slog.Error("error, error occured while initializing configuraion", "service_name", "groupie_tracker")
	}

	// Initializing and customization the slog logger and them set this logger as default logger
	_ = logger.Init(cfg)

	// app.Run() ...

	//
}
