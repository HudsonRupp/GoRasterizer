package main

import (
	"image"
	"math"
)

type Rasterizer struct {
	Width  int
	Height int
	ZBuf   []float64

	//TODO: consider separate struct for vertex-related buffers
	viewBuf   []Vec3
	rasterBuf []Vec3
	uvBuf     []Vec2

	blankFrame []byte
}

func NewRasterizer(width int, height int) *Rasterizer {

	zbuf := make([]float64, width*height)
	blank := make([]byte, width*height*4)
	for i := 0; i < len(blank); i += 4 {
		blank[i] = 100
		blank[i+1] = 100
		blank[i+2] = 100
		blank[i+3] = 255
	}

	return &Rasterizer{
		Width:      width,
		Height:     height,
		ZBuf:       zbuf,
		viewBuf:    make([]Vec3, 0, 1000),
		rasterBuf:  make([]Vec3, 1000),
		uvBuf:      make([]Vec2, 1000),
		blankFrame: blank,
	}
}

func (r *Rasterizer) ensureBufSize(needed int) {
	if cap(r.viewBuf) < needed {
		r.viewBuf = make([]Vec3, needed)
		r.rasterBuf = make([]Vec3, needed)
		r.uvBuf = make([]Vec2, needed)
	}
	r.viewBuf = r.viewBuf[:needed]
	r.rasterBuf = r.rasterBuf[:needed]
	r.uvBuf = r.uvBuf[:needed]
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

func (r *Rasterizer) Render(frame *image.RGBA, cam *Camera, meshes []*Mesh) {

	// Reset frame
	copy(frame.Pix, r.blankFrame)

	// Reset depth buffer
	for i := 0; i < len(r.ZBuf); i++ {
		r.ZBuf[i] = math.Inf(1)
	}

	viewMat := cam.getCameraMatrix()

	for _, mesh := range meshes {
		r.ensureBufSize(len(mesh.Vertices))

		for i, v := range mesh.Vertices {
			r.viewBuf[i] = viewMat.MultVec3(v.Position)
			r.rasterBuf[i] = r.worldToRaster(cam, v.Position)
			r.uvBuf[i] = v.UV
		}

		for _, face := range mesh.Faces {
			// z culling
			if r.viewBuf[face[0]].Z >= 0 || r.viewBuf[face[1]].Z >= 0 || r.viewBuf[face[2]].Z >= 0 {
				continue
			}

			p0 := r.rasterBuf[face[0]]
			p1 := r.rasterBuf[face[1]]
			p2 := r.rasterBuf[face[2]]

			// back face culling - based on winding relative to camera in raster space
			area := r.edgeFunction(p0, p1, p2)
			if area <= 0 {
				continue
			}

			// view space depths for zbuf
			z0 := -r.viewBuf[face[0]].Z
			z1 := -r.viewBuf[face[1]].Z
			z2 := -r.viewBuf[face[2]].Z

			//UV
			uv0 := r.uvBuf[face[0]]
			uv1 := r.uvBuf[face[1]]
			uv2 := r.uvBuf[face[2]]

			// Test colors to better see individual triangles
			/*
				c0 := Vec3{1, 0, 0}
				c1 := Vec3{0, 1, 0}
				c2 := Vec3{0, 0, 1}
			*/

			// Bounding box
			maxY := int(math.Ceil(math.Max(p0.Y, math.Max(p1.Y, p2.Y))))
			minY := int(math.Floor(math.Min(p0.Y, math.Min(p1.Y, p2.Y))))
			maxX := int(math.Ceil(math.Max(p0.X, math.Max(p1.X, p2.X))))
			minX := int(math.Floor(math.Min(p0.X, math.Min(p1.X, p2.X))))

			if minX < 0 {
				minX = 0
			}
			if minY < 0 {
				minY = 0
			}
			if maxX > frame.Rect.Dx() {
				maxX = frame.Rect.Dx()
			}
			if maxY > frame.Rect.Dy() {
				maxY = frame.Rect.Dy()
			}

			for y := int(minY); y < int(maxY); y++ {
				for x := int(minX); x < int(maxX); x++ {
					p := Vec3{X: float64(x) + .5, Y: float64(y) + .5} // middle of pixel
					inTriangle, w0, w1, w2 := r.pointInTriangle(p0, p1, p2, p)
					if inTriangle {
						// Perspective-correct Z
						inverseZ := w0*(1.0/z0) + w1*(1.0/z1) + w2*(1.0/z2)
						pixelZ := 1.0 / inverseZ

						// Perspective-correct UV interpolation TODO: review this
						uvX := pixelZ * (w0*(uv0.X/z0) + w1*(uv1.X/z1) + w2*(uv2.X/z2))
						uvY := pixelZ * (w0*(uv0.Y/z0) + w1*(uv1.Y/z1) + w2*(uv2.Y/z2))
						uv := Vec2{uvX, uvY}

						if pixelZ < r.ZBuf[y*r.Width+x] {
							r.ZBuf[y*r.Width+x] = pixelZ

							/*
								red := w0*c0.X + w1*c1.X + w2*c2.X
								green := w0*c0.Y + w1*c1.Y + w2*c2.Y
								blue := w0*c0.Z + w1*c1.Z + w2*c2.Z
							*/
							// sampleTex(uv.X, uv.Y, ...)

							r.setPixel(frame, x, y, SampleTexture(uv.X, uv.Y))
						}
					}
				}
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
