# GoRasterizer

3D Rasterizer from scratch in Go.

## Execution 

If on windows:
```bash
export GOOS=windows
```

To run:
```bash
go run . obj/*scene_name*
```

To import custom objects:
1. Make a new folder in the obj folder
2. Create three subfolders: mesh, skybox, texture
3. Put your .obj files in the mesh folder, the textures for the .obj files in the texture folder with the same name (ex. object1.obj -> object1.png), and skybox texture in skybox
4. Run with the path to the folder you created

## Controls
- WASD to move camera forward/back + left/right, Q/E to move camera up/down 
- Use arrow keys to change the direction of the camera

## Preview
<img width="828" height="716" alt="image" src="https://github.com/user-attachments/assets/22a59233-8e00-4fdc-affd-bf1b41201f0d" />
