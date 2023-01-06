package main

import (
	"fmt"

	"github.com/sf1/go-card/smartcard"
)

func main() {

	// establish reader context
	context, err := smartcard.EstablishContext()
	if err != nil {
		fmt.Println("Error establishing context")
		return
	}

	defer context.Release()
}
