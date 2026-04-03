package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Parse_OBJ(filename string) (*Mesh, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mesh := &Mesh{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "g") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		prefix := fields[0]

		if prefix == "v" {
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			mesh.Vertices = append(mesh.Vertices, Vec3{X: x, Y: y, Z: z})
		} else if prefix == "f" {
			var faceIndices []int

			for _, fStr := range fields[1:] {
				parts := strings.Split(fStr, "/")
				// Just want vertex for now, skip /vt/vn
				idx, err := strconv.Atoi(parts[0])
				if err == nil {
					faceIndices = append(faceIndices, idx-1)
				}
			}

			// Triangle Fan for > 3 vertices
			for i := 1; i < len(faceIndices)-1; i++ {
				tri := [3]int{
					faceIndices[0],
					faceIndices[i],
					faceIndices[i+1],
				}
				mesh.Faces = append(mesh.Faces, tri)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mesh, nil

}
