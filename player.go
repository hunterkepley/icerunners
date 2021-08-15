package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// PlayerState is a type for the current state of the Player
type PlayerState int

const (
	// PIdle ... PLAYERSTATE ENUM [1]
	PIdle PlayerState = iota + 1
	// PMovingLeft ... PLAYERSTATE ENUM [2]
	PMovingLeft
	// PMovingRight ... PLAYERSTATE ENUM [3]
	PMovingRight
	// PJumping ... PLAYERSTATE ENUM [4]
	PJumping
	// PFalling ... PLAYERSTATE ENUM [5]
	PFalling
)

func (b PlayerState) String() string {
	return [...]string{"Unknown",
		"PIdle", "PMovingLeft", "PMovingRight", "PJumping", "PFalling"}[b]
}

// PlayerType is a type for a local or network player (controllable)
type PlayerType int

const (
	PLocal PlayerType = iota + 1
	PNetwork
)

func (b PlayerType) String() string {
	return [...]string{"Undefined", "PLocal", "PNetwork"}[b]
}

// Player... The player!
type Player struct {
	position Vec2f
	size     Vec2i

	velocity  Vec2f
	moveSpeed float64

	state PlayerState
	_type PlayerType

	rotation float64

	subImageRect image.Rectangle
	image        *ebiten.Image
}

func createPlayer(position Vec2f, _type PlayerType) Player {
	return Player{
		position:     position,
		moveSpeed:    0.15,
		_type:        _type,
		subImageRect: image.Rect(0, 0, iPlayerSpritesheet.Bounds().Dx(), iPlayerSpritesheet.Bounds().Dy()),
		image:        iPlayerSpritesheet,
	}
}

func (p *Player) update() {
	p.size = newVec2i(p.subImageRect.Dx(), p.subImageRect.Dy())
	p.input()
	// Idle -- stop moving
	if p.state == PIdle {
		if p.velocity.x > 0 {
			p.velocity.x -= p.moveSpeed * 3
		} else if p.velocity.x < 0 {
			p.velocity.x += p.moveSpeed * 3
		}
		if p.velocity.x < p.moveSpeed && p.velocity.x > -p.moveSpeed {
			p.velocity.x = 0
		}
	}
	p.position.x += p.velocity.x
	p.position.y += p.velocity.y
}

func (p *Player) input() {
	p.state = PIdle

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.velocity.x -= p.moveSpeed
		p.state = PMovingLeft
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.velocity.x += p.moveSpeed
		p.state = PMovingRight
	}
}

func (p *Player) render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// ROTATE & FLIP
	op.GeoM.Translate(float64(0-p.size.x)/2, float64(0-p.size.y)/2)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(float64(p.size.x)/2, float64(p.size.y)/2)
	/*
		--- FOR ANIMATION ---
		p.subImageRect = image.Rect(
			p.spritesheet.sprites[p.animation.currentFrame].startPosition.x,
			p.spritesheet.sprites[p.animation.currentFrame].startPosition.y,
			p.spritesheet.sprites[p.animation.currentFrame].endPosition.x,
			p.spritesheet.sprites[p.animation.currentFrame].endPosition.y,
		)
	*/
	// POSITION
	op.GeoM.Translate(float64(p.position.x), float64(p.position.y))
	op.Filter = ebiten.FilterNearest // Maybe fix rotation grossness?

	// Render walk smoke
	/*
		--- PROBABLY A 'LATER' FEATURE, FROM UNRAY ---
		p.walkSmokeEmitter.render(screen)
	*/

	screen.DrawImage(p.image.SubImage(p.subImageRect).(*ebiten.Image), op)
}
