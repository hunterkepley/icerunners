package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

const (
	screenMultiplier = 1
	screenWidth      = 1366 * screenMultiplier
	screenHeight     = 768 * screenMultiplier
)

var (
	mdataFont    font.Face // For fonts later
	mversionFont font.Face

	gameInitialized bool
)

// Game is the info for the game
type Game struct {
	gameEntities GameEntities

	state int // The game state, 0 is in main menu, 1 is in game, 2 is paused

	settings Settings // Game settings

}

// Init initializes the game
func (g *Game) Init() {

	g.gameEntities.init()

	// Fonts
	g.InitFonts()

	// State starts in game [temporary]
	g.state = 1

	// Init music
	//loadMusic()
	// Play song
	//go music[rand.Intn(len(music)-1)].play()

	// GAME SETTINGS
	loadSettings(&g.settings)

	if g.settings.Graphics.Fullscreen { // Enable fullscreen if enabled
		ebiten.SetFullscreen(true)
	}

}

// Update updates the game
func (g *Game) Update(screen *ebiten.Image) error {
	if !gameInitialized {
		g.Init()
		gameInitialized = true
	}

	// Update game
	if g.state == 1 {
		updateGame(&g.gameEntities)
	}

	return nil
}

// Draw renders everything!
func (g *Game) Draw(screen *ebiten.Image) {

	// Draw game
	if g.state == 1 { // inGame
		drawGame(&g.gameEntities, screen)
	}

}

// Layout is the screen layout?...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	loadPregameResources()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle(fmt.Sprintf("I C E R U N N E R S %s", gameVersionType))
	ebiten.SetWindowResizable(true)

	// Load icon
	//icon16x16, _ := loadRegularImage("./Assets/Art/Icon/icon16.png")
	//icon32x32, _ := loadRegularImage("./Assets/Art/Icon/icon32.png")
	//icon48x48, _ := loadRegularImage("./Assets/Art/Icon/icon48.png")
	//icon64x64, _ := loadRegularImage("./Assets/Art/Icon/icon64.png")
	//ebiten.SetWindowIcon([]image.Image{icon16x16, icon32x32, icon48x48, icon64x64})

	// Hide cursor
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func loadPregameResources() {
	// Images
	loadTileImages()
	//loadUIImages()
}
