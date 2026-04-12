package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
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

	var sceneName string
	if len(os.Args) > 1 {
		sceneName = os.Args[1]
	} else {
		sceneName = "obj/blank"
		folders := []string{"mesh", "texture", "skybox"}
		for _, folder := range folders {
			dirPath := filepath.Join(sceneName, folder)

			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				log.Printf("Error creating blank scene folders: %v", err)
				continue
			}
		}
	}

	meshDir := filepath.Join(sceneName, "mesh")
	texDir := filepath.Join(sceneName, "texture")
	skyboxDir := filepath.Join(sceneName, "skybox")

	var meshes []*Mesh

	skyMesh, err := LoadOBJ("obj/skysphere.obj")
	if err != nil {
		log.Printf("Failed to load sky mesh: %v", err)
	}
	skyMesh.IsSky = true

	skyboxFiles, err := os.ReadDir(skyboxDir)
	if err != nil {
		log.Printf("Failed to read skybox dir: %v", err)
	}
	if len(skyboxFiles) > 0 {
		for _, file := range skyboxFiles {
			if !file.IsDir() {
				skyTexPath := filepath.Join(skyboxDir, file.Name())
				skyTex, err := LoadTexture(skyTexPath)
				if err != nil {
					log.Printf("Failed to load sky texture %s: %v", skyTexPath, err)
				}

				skyMesh.Texture = skyTex

				log.Printf("Loaded Sky")
				break
			}
		}
	}

	meshes = append(meshes, skyMesh)

	meshFiles, err := os.ReadDir(meshDir)
	if err != nil {
		log.Fatalf("Failed to read mesh dir")
	}

	for _, file := range meshFiles {
		if file.IsDir() || filepath.Ext(file.Name()) != ".obj" {
			continue
		}

		objPath := filepath.Join(meshDir, file.Name())
		mesh, err := LoadOBJ(objPath)
		if err != nil {
			log.Printf("Failed to load obj %s: %v", objPath, err)
			continue
		}

		baseName := strings.TrimSuffix(file.Name(), ".obj")
		texPath := filepath.Join(texDir, baseName+".png")

		_, err = os.Stat(texPath)

		if err == nil {
			tex, err := LoadTexture(texPath)
			if err != nil {
				log.Printf("Failed to load texture %s: %v", texPath, err)
			} else {
				mesh.Texture = tex
			}
		} else {
			log.Printf("No matching texture found for %s", texPath)
		}

		meshes = append(meshes, mesh)
	}

	if len(meshFiles) == 0 {
		meshes = append(meshes, TestQuad())
	}

	engine := NewRasterizer(1920, 1080)

	game := NewGame(engine, 1920, 1080, meshes)

	Run(game)
}
