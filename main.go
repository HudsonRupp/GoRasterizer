package main

import (
	"log"
	"os"
)

func main() {
	f, _ := os.Create("debug.log")
	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting...")

	log.Println(os.Args)
	if len(os.Args) > 1 {
		Parse_OBJ(os.Args[1])
	}

	engine := NewRasterizer(720, 720)

	game := NewGame(engine, 720, 720)

	Run(game)
}
