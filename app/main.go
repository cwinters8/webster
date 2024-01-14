package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func setup() error {
	run := app.New()
	window := run.NewWindow("Webster")
	logo := canvas.NewImageFromFile("./assets/images/webster.png")
	run.SetIcon(logo.Resource)
	window.SetContent(widget.NewLabel("Hello, World!"))
	window.ShowAndRun()
	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
}
