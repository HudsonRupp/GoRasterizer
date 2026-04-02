package main

import (
	"image"
	"math"
)

type Rasterizer struct {
	Width  int
	Height int
	ZBuf   [][]int
}

func NewRasterizer(width int, height int) *Rasterizer {

	zbuf := make([][]int, width)
	for i := range zbuf {
		zbuf[i] = make([]int, height)
	}

	return &Rasterizer{
		Width:  width,
		Height: height,
		ZBuf:   zbuf,
	}
}

func rotate(v Vec3, angle float64) Vec3 {
	rad := angle * (math.Pi / 180.0)
	oZ := -5.0
	rZ := v.Z - oZ

	nX := v.X*math.Cos(rad) + rZ*math.Sin(rad)
	nZ := -v.X*math.Sin(rad) + rZ*math.Sin(rad)

	return Vec3{
		X: float64(nX),
		Y: v.Y,
		Z: float64(nZ + oZ),
	}
}

var angle float64 = 0

func (r *Rasterizer) Render(frame *image.RGBA, cam *Camera) {

	// Set to black
	for i := 0; i < len(frame.Pix); i += 4 {
		frame.Pix[i] = 0     //R
		frame.Pix[i+1] = 0   //G
		frame.Pix[i+2] = 0   //B
		frame.Pix[i+3] = 255 // Alpha
	}

	angle += 1
	if angle == 90 {
		angle = 270
	} else if angle == 360 {
		angle = 0
	}

	V0 := Vec3{0, 1.5, -5}
	V1 := Vec3{-1.5, -1.5, -5}
	V2 := Vec3{1.5, -1.5, -5}
	c0 := Vec3{1, 0, 0}
	c1 := Vec3{0, 1, 0}
	c2 := Vec3{0, 0, 1}

	//V0 = rotate(V0, angle)
	//V1 = rotate(V1, angle)
	//V2 = rotate(V2, angle)

	p0 := r.worldToRaster(cam, V0)
	p1 := r.worldToRaster(cam, V1)
	p2 := r.worldToRaster(cam, V2)

	// maxY := math.Max(p0.Y, math.Max(p1.Y, p2.Y))
	// minY := math.Min(p0.Y, math.Min(p1.Y, p2.Y))
	// maxX := math.Max(p0.X, math.Max(p1.X, p2.X))
	// minX := math.Min(p0.X, math.Min(p1.X, p2.X))

	// y := 0; y < frame.Rect.Dy(); y++
	//for y := int(minY); y < int(maxY); y++ {
	for y := 0; y < frame.Rect.Dy(); y++ {
		//for x := int(minX); x < int(maxX); x++ {
		for x := 0; x < frame.Rect.Dx(); x++ {
			p := Vec3{X: float64(x) + .5, Y: float64(y) + .5} // middle of pixel
			inTriangle, w0, w1, w2 := r.pointInTriangle(p0, p1, p2, p)
			if inTriangle {
				red := w0*c0.X + w1*c1.X + w2*c2.X
				green := w0*c0.Y + w1*c1.Y + w2*c2.Y
				blue := w0*c0.Z + w1*c1.Z + w2*c2.Z
				r.setPixel(frame, x, y, Vec3{red, green, blue})
			}
		}
	}
}

func (r *Rasterizer) setPixel(frame *image.RGBA, x int, y int, color Vec3) {
	if x < 0 || y < 0 || x >= frame.Rect.Dx() || y >= frame.Rect.Dy() {
		return
	}

	i := y*frame.Stride + x*4

	frame.Pix[i] = uint8(color.X * 255)
	frame.Pix[i+1] = uint8(color.Y * 255)
	frame.Pix[i+2] = uint8(color.Z * 255)
	frame.Pix[i+3] = 255 //always opaque
}

func getProjMatrix(fov, near, far float64) Mat4 {
	m := Mat4{}
	scale := 1 / math.Tan(fov*0.5*math.Pi/180)
	m[0][0] = scale                      // scale x
	m[1][1] = scale                      // scale y
	m[2][2] = -far / (far - near)        // remap z [0,1]
	m[3][2] = -far * near / (far - near) // remap z [0, 1]
	m[2][3] = -1                         // w = -z

	return m
}

func (r *Rasterizer) worldToRaster(c *Camera, pWorld Vec3) Vec3 {
	pNDC := c.VPMatrix.MultVec3(pWorld)
	return Vec3{X: ((pNDC.X + 1) * 0.5 * float64(r.Width)), Y: ((1 - (pNDC.Y+1)*0.5) * float64(r.Height)), Z: -pNDC.Z}
}

func (ras *Rasterizer) projectToRaster(pCam Vec3, width int, height int) Vec3 {
	nearClippingPlane := 0.1
	aspectRatio := 1.0
	fov := 90 * (math.Pi / 180.0)

	pScreen := Vec3{}
	pScreen.X = nearClippingPlane * pCam.X / -pCam.Z
	pScreen.Y = nearClippingPlane * pCam.Y / -pCam.Z

	t := math.Tan(fov/2) * nearClippingPlane
	r := t * aspectRatio
	b := -t
	l := -r

	pNDC := Vec3{}
	pNDC.X = 2*pScreen.X/(r-l) - (r+l)/(r-l)
	pNDC.Y = 2*pScreen.Y/(t-b) - (t+b)/(t-b)

	pRas := Vec3{}
	pRas.X = (pNDC.X + 1) / 2 * float64(width)
	pRas.Y = (1 - pNDC.Y) / 2 * float64(height)
	pRas.Z = -pCam.Z

	//log.Printf("pCam: %v, pScreen: %v, pNDC: %v, pRas: %v", pCam, pScreen, pNDC, pRas)
	return pRas
}

func (ras *Rasterizer) edgeFunction(v1, v2, p Vec3) float64 {
	return ((p.X-v1.X)*(v2.Y-v1.Y) - (p.Y-v1.Y)*(v2.X-v1.X))
}

func (r *Rasterizer) pointInTriangle(V0, V1, V2, p Vec3) (bool, float64, float64, float64) {
	area := r.edgeFunction(V0, V1, V2)
	w0 := r.edgeFunction(V0, V1, p)
	w1 := r.edgeFunction(V1, V2, p)
	w2 := r.edgeFunction(V2, V0, p)

	if w0 >= 0 && w1 >= 0 && w2 >= 0 {
		w0 /= area
		w1 /= area
		w2 /= area
		return true, w0, w1, w2
	}
	return false, w0, w1, w2
}
