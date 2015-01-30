package main

import (
	"os"

	"github.com/frenata/gomua"
)

func main() {
	// write a message on the stdin and send it
	gomua.Send(gomua.WriteMessage(os.Stdin, os.Stdin))
}
