package main

import (
	"log"
	"os"
	//"runtime/pprof"
)

func main() {
	/*
		pf, _ := os.Create("cpu.prof")
		defer pf.Close()
		pprof.StartCPUProfile(pf)
		defer pprof.StopCPUProfile()
	*/

	f, _ := os.Create("debug.log")
	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting...")

	log.Println(os.Args)
	var mesh *Mesh
	var err error
	if len(os.Args) > 1 {
		mesh, err = Parse_OBJ(os.Args[1])
	} else {
		mesh = TestQuad()
	}

	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	engine := NewRasterizer(1920, 1080)

	game := NewGame(engine, 1920, 1080, []*Mesh{mesh})

	Run(game)
}
