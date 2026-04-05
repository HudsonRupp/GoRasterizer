package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Renderer interface {
	Render(frame *image.RGBA)
}

type Game struct {
	rasterizer *Rasterizer
	camera     *Camera
	buffer     *image.RGBA
	offscreen  *ebiten.Image
	meshes     []*Mesh
}

func NewGame(r *Rasterizer, width int, height int, meshes []*Mesh) *Game {
	return &Game{
		rasterizer: r,
		camera:     NewCamera(85, 0.1, 100, width, height),
		buffer:     image.NewRGBA(image.Rect(0, 0, width, height)),
		offscreen:  ebiten.NewImage(width, height),
		meshes:     meshes,
	}
}

func (g *Game) Update() error {
	moved := false
	translateSpeed := 0.1

	rotateSpeed := 2.0 // deg/frame

	//Rotation
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.Rotate(-rotateSpeed, 0)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.Rotate(rotateSpeed, 0)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Rotate(0, rotateSpeed)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Rotate(0, -rotateSpeed)
		moved = true
	}

	forward := g.camera.Target.Sub(g.camera.Position).Normalized()
	right := (Vec3{X: 0, Y: 1, Z: 0}).Cross(forward).Normalized()
	dForward := Vec3{forward.X * translateSpeed, forward.Y * translateSpeed, forward.Z * translateSpeed}
	dRight := Vec3{right.X * translateSpeed, right.Y * translateSpeed, right.Z * translateSpeed}

	//Translation
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camera.Translate(dForward)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camera.Translate(dRight)
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camera.Translate(dForward.ScalarMult(-1))
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camera.Translate(dRight.ScalarMult(-1))
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.Translate(Vec3{0, translateSpeed, 0})
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.Translate(Vec3{0, -translateSpeed, 0})
		moved = true
	}

	if moved {
		g.camera.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.rasterizer.Render(g.buffer, g.camera, g.meshes)
	g.offscreen.WritePixels(g.buffer.Pix) // Copy buffer to ebiten.Image
	screen.DrawImage(g.offscreen, nil)
	m := fmt.Sprintf("TPS: %.2f, FPS: %.2f, CAM: %v, LOOKAT: %v", ebiten.ActualTPS(), ebiten.ActualFPS(), g.camera.Position, g.camera.Target)
	ebitenutil.DebugPrint(screen, m)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Width: %v, Height: %v", g.camera.Width, g.camera.Height), 0, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.camera.Width, g.camera.Height
}

func Run(game *Game) {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("GoRasterizer")
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
