package main

import (
	"log"

	"github.com/mgironi/operation-fire-quasar/location"
	"github.com/mgironi/operation-fire-quasar/message"
	"github.com/mgironi/operation-fire-quasar/store"
)

func main() {
	log.SetFlags(0)

	// checks and console display, if only asked for help menu/instructions
	if AskForHelp() {
		return
	}

	// initialices the store (in memory)
	store.Initialize()

	// parse
	distances, messages, parseErr := ParseArgs()
	if parseErr != nil {
		log.Fatalf("ERROR\t%s", parseErr.Error())
	}

	// Gets location
	x, y := GetLocation(distances...)
	log.Printf("The location coordinates is x: %f, y: %f", x, y)

	// Gets complete message
	message := GetMessage(messages...)
	log.Printf("The complete message is '%s'.", message)
}

// input: distance to the transmitter recieved on each satlelite
// output: the coordinates 'x' and 'y' of the message emiter
func GetLocation(distances ...float32) (x, y float32) {
	return location.GetLocation(distances...)
}

// input: the message as it is recieved on each satelite
// output: the message as it is generated by the transmitter
func GetMessage(messages ...[]string) (msg string) {
	return message.GetMessage(messages...)
}
