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

	var rawPositions []Vec3
	var rawUVs []Vec2
	var rawNormals []Vec3

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

		switch fields[0] {
		case "v":
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			rawPositions = append(rawPositions, Vec3{x, y, z})
		case "vt":
			u, _ := strconv.ParseFloat(fields[1], 64)
			v, _ := strconv.ParseFloat(fields[2], 64)
			rawUVs = append(rawUVs, Vec2{u, v})
		case "vn":
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			rawNormals = append(rawNormals, Vec3{x, y, z})
		case "f":
			var faceIndices []int

			for _, fStr := range fields[1:] {
				parts := strings.Split(fStr, "/")

				vInd, _ := strconv.Atoi(parts[0])
				uvInd := 0
				nInd := 0

				if len(parts) > 1 && parts[1] != "" {
					uvInd, _ = strconv.Atoi(parts[1])
				}
				if len(parts) > 2 && parts[2] != "" {
					nInd, _ = strconv.Atoi(parts[2])
				}

				vertex := Vertex{
					Position: rawPositions[vInd-1],
				}
				if uvInd > 0 {
					vertex.UV = rawUVs[uvInd-1]
				}
				if nInd > 0 {
					vertex.Normal = rawNormals[nInd-1]
				}

				mesh.Vertices = append(mesh.Vertices, vertex)

				faceIndex := len(mesh.Vertices) - 1 // last added vertex
				faceIndices = append(faceIndices, faceIndex)
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

func SampleTexture(u float64, v float64) (color Vec3) {
	//checkerboard for now, will read from texture image in future
	scale := 8.0
	if (int(u*scale)+int(v*scale))%2 == 1 {
		return Vec3{1, 1, 1}
	}
	return Vec3{0, 0, 0}

}
