package main

import (
	"log"
	"os"

	"github.com/Pikita-noname/davinchikTgApp/internal/view"
)

func main() {
	app := view.NewApp()

	if err := app.View.SetRoot(app.MainMenu.View, true).Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
		os.Exit(1)
	}
}
