package main

import (
	gui "GoGUI/src/GUI"
	"log"
)

func main() {
	params, err := gui.GetParams()
	if err != nil {
		log.Fatal(err)
	}

	gui.Graph(params)
}
