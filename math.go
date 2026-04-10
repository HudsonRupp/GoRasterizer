package main

import "math"

type Mat4 [4][4]float64
type Mat2 [2][2]float64

func (m Mat4) MultVec3(p Vec3) Vec3 {

	x := p.X*m[0][0] + p.Y*m[1][0] + p.Z*m[2][0] + m[3][0]
	y := p.X*m[0][1] + p.Y*m[1][1] + p.Z*m[2][1] + m[3][1]
	z := p.X*m[0][2] + p.Y*m[1][2] + p.Z*m[2][2] + m[3][2]
	w := p.X*m[0][3] + p.Y*m[1][3] + p.Z*m[2][3] + m[3][3]

	if w != 1 && w != 0 {
		return Vec3{X: x / w, Y: y / w, Z: z / w}
	}
	return Vec3{X: x, Y: y, Z: z}
}

func (m Mat4) MultMat4(o Mat4) Mat4 {
	var r Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			r[i][j] = dot(m, o, i, j)
		}
	}
	return r
}

func (m Mat4) Transpose() Mat4 {
	var r Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			r[i][j] = m[j][i]
		}
	}
	return r
}

func dot(a, b Mat4, row, col int) float64 {
	return a[row][0]*b[0][col] +
		a[row][1]*b[1][col] +
		a[row][2]*b[2][col] +
		a[row][3]*b[3][col]
}

func (m Mat2) det() float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Normalized() Vec3 {
	m := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if m == 0 {
		return Vec3{X: 0, Y: 0, Z: 0}
	}
	return Vec3{X: v.X / m, Y: v.Y / m, Z: v.Z / m}
}

func (v Vec3) Dot(o Vec3) float64 {
	return v.X*o.X + v.Y + o.Y*v.Z*o.Z
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X}
}

func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3{X: v.X - o.X, Y: v.Y - o.Y, Z: v.Z - o.Z}
}

func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3{X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z}
}

func (v Vec3) ScalarMult(k float64) Vec3 {
	return Vec3{X: v.X * k, Y: v.Y * k, Z: v.Z * k}
}

func FaceNormal(a Vec3, b Vec3, c Vec3) Vec3 {
	return (b.Sub(a)).Cross((c.Sub(a))).Normalized()
}

type Vec2 struct {
	X, Y float64
}
