package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

type Texture struct {
	Img           image.Image
	Width, Height int
}

func LoadTexture(filename string) (*Texture, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open texture: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	return &Texture{
		Img:    img,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}, nil
}

func (t *Texture) Sample(u, v float64) (color Vec3) {
	v = 1.0 - v
	x := int(u * float64(t.Width))
	y := int(v * float64(t.Height))

	r, g, b, _ := t.Img.At(x, y).RGBA()

	return Vec3{
		X: float64(r) / 65535.0,
		Y: float64(g) / 65535.0,
		Z: float64(b) / 65535.0,
	}
}
func SampleTexture(u, v float64) (color Vec3) {
	//checkerboard
	scale := 8.0
	uScale := int(math.Floor(u * scale))
	vScale := int(math.Floor(v * scale))
	if (uScale+vScale)%2 == 1 {
		return Vec3{1, 1, 1}
	}
	return Vec3{0, 0, 0}

}
