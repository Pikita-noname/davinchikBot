package main

import (
	"log"
	"os"

	"github.com/Pikita-noname/davinchikTgApp/internal/view"
)

func main() {
	app := view.NewApp()

	// Убедимся, что MainMenu и Options созданы перед вызовом SetRoot
	if app.MainMenu.View == nil {
		log.Fatal("MainMenu.View is nil")
	}
	if app.Options.View == nil {
		log.Fatal("Options.View is nil")
	}

	if err := app.View.SetRoot(app.MainMenu.View, true).Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
		os.Exit(1)
	}
}
