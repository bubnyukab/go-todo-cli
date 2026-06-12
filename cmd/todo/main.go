package main

import (
	"log"

	"github.com/bubnyukab/go-todo-cli/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("unable to run the app: %v", err)
	}
}
