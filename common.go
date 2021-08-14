package main

// ^ DIRECTION ENUM ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	// Pi is pi
	Pi = 3.14159265358979323846264338327950288419716939937510582097494459 // https://oeis.org/A000796
)

// Vec2f is a vector of 2 float64's
type Vec2f struct {
	x float64
	y float64
}

func newVec2f(x float64, y float64) Vec2f {
	return Vec2f{x, y}
}

// Vec2i is a vector of 2 int's
type Vec2i struct {
	x int
	y int
}

func newVec2i(x int, y int) Vec2i {
	return Vec2i{x, y}
}

// Vec2b is a vector of 2 bool's
type Vec2b struct {
	x bool
	y bool
}

func newVec2b(x bool, y bool) Vec2b {
	return Vec2b{x, y}
}

func vec2f(x Vec2i) Vec2f {
	return newVec2f(float64(x.x), float64(x.y))
}

func vec2i(x Vec2f) Vec2i {
	return newVec2i(int(x.x), int(x.y))
}

// Rolumn is a row/column vector
type Rolumn struct {
	row    int
	column int
}

func newRolumn(row int, column int) Rolumn {
	return Rolumn{row, column}
}
