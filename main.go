package main

import (
	"log"
)

func main() {
	dji, err := NewDjiModem()
	if err != nil {
		log.Fatal(err)
	}

	dji.AtShell()
}
