package main

import "math"

type Camera struct {
	Position               Vec3
	Target                 Vec3
	FOV, NearClip, FarClip float64
	Width, Height          int

	VPMatrix Mat4
}

func NewCamera(fov, nearClip, farClip float64, width, height int) *Camera {
	cam := &Camera{
		Position: Vec3{X: 0, Y: 0, Z: 0},
		Target:   Vec3{X: 0, Y: 0, Z: -5},
		FOV:      fov,
		NearClip: nearClip,
		FarClip:  farClip,
		Width:    width,
		Height:   height,
	}
	cam.Update()
	return cam
}

func (c *Camera) Translate(d Vec3) {
	c.Position = c.Position.Add(d)
	c.Target = c.Target.Add(d)
}

func (c *Camera) Update() {
	viewMatrix := c.getCameraMatrix()
	projMatrix := c.getProjMatrix()

	c.VPMatrix = viewMatrix.MultMat4(projMatrix)
}

func (c *Camera) getCameraMatrix() Mat4 {
	forward := c.Position.Sub(c.Target).Normalized()
	tmp := Vec3{X: 0, Y: 1, Z: 0}
	right := tmp.Cross(forward).Normalized()
	up := forward.Cross(right)

	tx := -(c.Position.X*right.X + c.Position.Y*right.Y + c.Position.Z*right.Z)
	ty := -(c.Position.X*up.X + c.Position.Y*up.Y + c.Position.Z*up.Z)
	tz := -(c.Position.X*forward.X + c.Position.Y*forward.Y + c.Position.Z*forward.Z)
	return Mat4{
		{right.X, right.Y, right.Z, 0},
		{up.X, up.Y, up.Z, 0},
		{forward.X, forward.Y, forward.Z, 0},
		{tx, ty, tz, 1},
	}
}

func (c *Camera) getProjMatrix() Mat4 {
	m := Mat4{}
	aspectRatio := float64(c.Width) / float64(c.Height)
	top := math.Tan(c.FOV/2) * c.NearClip
	bottom := -top
	right := top * aspectRatio
	left := bottom
	m[0][0] = (2 * c.NearClip) / (right - left)
	m[1][1] = (2 * c.NearClip) / (top - bottom)
	m[2][0] = (right + left) / (right - left)
	m[2][1] = (top + bottom) / (top - bottom)
	m[2][2] = -(c.FarClip + c.NearClip) / (c.FarClip - c.NearClip)
	m[2][3] = -1
	m[3][2] = -(2 * c.FarClip * c.NearClip) / (c.FarClip - c.NearClip)

	/*scale := 1 / math.Tan(c.FOV*0.5*math.Pi/180)
	m[0][0] = scale                                              // scale x
	m[1][1] = scale                                              // scale y
	m[2][2] = -c.FarClip / (c.FarClip - c.NearClip)              // remap z [0,1]
	m[3][2] = -c.FarClip * c.NearClip / (c.FarClip - c.NearClip) // remap z [0, 1]
	m[2][3] = -1                                                 // w = -z
	*/
	return m
}
