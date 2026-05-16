package main

import (
	"context"
	"log"

	"github.com/kaustubha-chaturvedi/Rook/internal/execution"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	service := execution.NewService()

	app := application.New(application.Options{
		Name:        "Rook",
		Description: "Production-safe SQL notebook",
		Services: []application.Service{
			application.NewService(service),
		},
		OnShutdown: func() {
			service.Shutdown(context.Background())
		},
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
