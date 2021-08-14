package main

import "github.com/hajimehoshi/ebiten"

var (
	// i prefix is for images
	iTileSpritesheet   *ebiten.Image
	iPlayerSpritesheet *ebiten.Image
	//iUISpritesheet             *ebiten.Image
)

/*
func loadUIImages() {
	iUISpritesheet, _ = loadImage("./Assets/UI/spritesheet.png")
}*/

func loadTileImages() {
	iTileSpritesheet, _ = loadImage("./Assets/Map/Tiles/spritesheet.png")
	iPlayerSpritesheet, _ = loadImage("./Assets/Player/spritesheet.png")
}
