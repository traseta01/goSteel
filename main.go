package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sf1/go-card/smartcard"
)

var reader *smartcard.Reader

func main() {

	// establish reader context
	context, err := smartcard.EstablishContext()
	if err != nil {
		fmt.Println("Error establishing context")
		return
	}
	// get readers list
	readers_list, err := context.ListReadersWithCard()
	if err != nil {
		return
	}

	// handle list of readers & get card
	if len(readers_list) == 0 {
		fmt.Println("\nPlease insert smart card\n")
		return
	}

	// work with one reader for now
	if len(readers_list) == 1 {
		reader = readers_list[0]
	}

	// connect to reader
	card, err := reader.Connect()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%T\n------------------------------------\n", card)

	defer card.Disconnect()

	defer context.Release()

	// create new fyne app
	myApp := app.New()
	windw := myApp.NewWindow("GoSteel")

	buttonScan := widget.NewButton("Scan", func() {

	})

	buttonGenPdf := widget.NewButton("Create PDF", func() {

	})

	contetButtons := container.New(layout.NewHBoxLayout(), buttonGenPdf, layout.NewSpacer(), buttonScan)

	// set fyne window contents and attributes
	windw.SetContent(container.New(layout.NewVBoxLayout(), contetButtons))
	windw.Resize(fyne.NewSize(800, 900))
	windw.SetFixedSize(true)
	windw.CenterOnScreen()

	windw.ShowAndRun()
}
