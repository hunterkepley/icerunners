package main

import (
	"github.com/hajimehoshi/ebiten"
)

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
	_, wY := ebiten.Wheel()
	if wY > 0 && c.zoom < 1.4 {
		c.zoom += 0.05
	} else if wY < 0 && c.zoom > .4 {
		c.zoom = ((c.zoom * 100) - 5) / 100
	}
	c.size = Vec2f{screenWidth + (screenWidth - screenWidth*c.zoom), screenHeight + (screenHeight - screenHeight*c.zoom)}
}
