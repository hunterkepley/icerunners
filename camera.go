package main

type Camera struct {
	position    Vec2f
	cameraSpeed float64
	zoom        float64
	size        Vec2f
}

func createCamera(position Vec2f, cameraSpeed float64) Camera {
	return Camera{
		position:    position,
		cameraSpeed: cameraSpeed,
		zoom:        1,
		size:        Vec2f{screenWidth, screenHeight},
	}
}

func (c *Camera) update() {
	c.size = Vec2f{screenWidth + (screenWidth - screenWidth*c.zoom), screenHeight + (screenHeight - screenHeight*c.zoom)}
}
