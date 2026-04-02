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
	keys       []ebiten.Key
}

func NewGame(r *Rasterizer, width int, height int) *Game {
	return &Game{
		rasterizer: r,
		camera:     NewCamera(90, 0.1, 100, width, height),
		buffer:     image.NewRGBA(image.Rect(0, 0, width, height)),
		offscreen:  ebiten.NewImage(width, height),
	}
}

func (g *Game) Update() error {
	moved := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camera.Translate(Vec3{0, 0, -0.1})
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camera.Translate(Vec3{-0.1, 0, 0})
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camera.Translate(Vec3{0, 0, 0.1})
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camera.Translate(Vec3{0.1, 0, 0})
		moved = true
	}

	if moved {
		g.camera.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.rasterizer.Render(g.buffer, g.camera)
	g.offscreen.WritePixels(g.buffer.Pix) // Copy buffer to ebiten.Image
	screen.DrawImage(g.offscreen, nil)    // Draw image
	m := fmt.Sprintf("TPS: %.2f, FPS: %.2f, CAM: %v, LOOKAT: %v", ebiten.ActualTPS(), ebiten.ActualFPS(), g.camera.Position, g.camera.Target)
	ebitenutil.DebugPrint(screen, m)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Width: %v, Height: %v", g.camera.Width, g.camera.Height), 0, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//return outsideWidth, outsideHeight
	return 720, 720
}

func Run(game *Game) {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("GoRasterizer")
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
