package main

import (
	"context"
	"log"

	"github.com/kaustubha-chaturvedi/Rook/internal/execution"
	"github.com/kaustubha-chaturvedi/Rook/internal/notebook"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	execService := execution.NewService()
	notebookService := notebook.NewService()

	app := application.New(application.Options{
		Name:        "Rook",
		Description: "Production-safe SQL notebook",
		Services: []application.Service{
			application.NewService(execService),
			application.NewService(notebookService),
		},
		OnShutdown: func() {
			execService.Shutdown(context.Background())
		},
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
