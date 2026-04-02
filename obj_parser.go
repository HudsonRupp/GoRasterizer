package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func Parse_OBJ(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	file_string := string(bytes)
	lines := strings.Split(file_string, "\n")

	var vertices []Vec3

	for i := 0; i < len(lines); i++ {
		log.Printf("%v", lines[i])
		if lines[i][0] == '#' || lines[i][0] == 'g' {
			// Comment or group (ignore FOR NOW)

		} else if lines[i][0] == 'f' {
			// face definition

		} else if lines[i][0] == 'v' {
			if lines[i][1] == 't' {
				// texture coordinates
			} else if lines[i][1] == 'n' {
				// normal definition

			} else {
				// Vertex definition
				coords := strings.Split(lines[i][2:], " ")
				x, _ := strconv.ParseFloat(coords[0], 64)
				y, _ := strconv.ParseFloat(coords[1], 64)
				z, _ := strconv.ParseFloat(coords[2], 64)
				vertices = append(vertices, Vec3{X: x, Y: y, Z: z})
			}
		}

	}

}
