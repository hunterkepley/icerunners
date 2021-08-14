package main

import (
	"github.com/hajimehoshi/ebiten"
)

// Map holds the image in which the world is drawn and moved (simulated camera)
type Map struct {
	cameraPosition  Vec2f
	cameraSize      Vec2f
	cameraMoveSpeed float64

	bgImage    *ebiten.Image
	image      *ebiten.Image
	imageScale float64
}

func initializeMap() Map {
	bgImage, _ := loadImage("./Assets/Map/bg.png")
	image, _ := loadImage("./Assets/Map/bg.png")

	return Map{
		cameraPosition:  Vec2f{0, 0},
		cameraMoveSpeed: 3.,

		bgImage: bgImage,
		image:   image,
	}
}

func (m *Map) update(camera Camera) {
	m.imageScale = camera.zoom
	m.cameraPosition = camera.position
	m.cameraSize = camera.size
}

func (m *Map) render(screen *ebiten.Image) {

	// Render map to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(m.cameraPosition.x-screenWidth/2, m.cameraPosition.y-screenHeight/2)
	op.GeoM.Scale(m.imageScale, m.imageScale)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(m.image, op)

}

func (m *Map) clearImage() {
	m.image.Clear()
	// Render map bg
	op := &ebiten.DrawImageOptions{}
	m.image.DrawImage(m.bgImage, op)
}
