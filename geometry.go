package main

type Mesh struct {
	Vertices []Vec3
	Faces    [][3]int
	Color    Vec3
	Position Vec3
}

func TestQuad() *Mesh {
	return &Mesh{
		Vertices: []Vec3{
			{-1, 1, 0},
			{-1, -1, 0},
			{1, -1, 0},
			{1, 1, 0},
		},
		Faces: [][3]int{
			{0, 1, 2},
			{0, 2, 3},
		},
		Color: Vec3{1, 0, 0},
	}
}
