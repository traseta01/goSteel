package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
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
		fmt.Println("\nplease insert smart card\n")
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

	myApp := app.New()
	windw := myApp.NewWindow("GoSteel")

	windw.ShowAndRun()
}
