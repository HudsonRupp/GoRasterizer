package main

type Mesh struct {
	Vertices []Vertex
	Faces    [][3]int
	Color    Vec3
	Position Vec3
}

type Vertex struct {
	Position Vec3
	UV       Vec2
	Normal   Vec3
}

func TestQuad() *Mesh {
	return &Mesh{
		Vertices: []Vertex{
			Vertex{
				Position: Vec3{-1, 1, 0},
				UV:       Vec2{-1, 1},
				Normal:   Vec3{0, 0, 1},
			},
			Vertex{
				Position: Vec3{-1, -1, 0},
				UV:       Vec2{-1, -1},
				Normal:   Vec3{0, 0, 1},
			},
			Vertex{
				Position: Vec3{1, -1, 0},
				UV:       Vec2{1, -1},
				Normal:   Vec3{0, 0, 1},
			},
			Vertex{
				Position: Vec3{1, 1, 0},
				UV:       Vec2{1, 1},
				Normal:   Vec3{0, 0, 1},
			},
		},
		Faces: [][3]int{
			{0, 1, 2},
			{0, 2, 3},
		},
		Color: Vec3{1, 0, 0},
	}
}
