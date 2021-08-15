package main

import (
	"fmt"
	"image"
	"time"

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

	//need states for: being stunned, hitting, invincible, and having frostrunners
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

	jump       bool
	playerType string
	stunned    bool
	iframe     bool
	hasBoot    bool

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

	if p.stunned == false {
		p.input()
	}

	if p.stunned == true {
		p.velocity.x = 0
	}
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

	//falling
	if p.position.y < 567 {
		p.velocity.y += p.moveSpeed
		p.state = PFalling
	}

	//on ground
	if p.position.y >= 568 {
		p.velocity.y = 0
		p.position.y = 567
		p.jump = false
	}

	p.position.x += p.velocity.x
	p.position.y += p.velocity.y

	//bounds
	if p.position.x < 0 {
		p.position.x = 0
	}

	if p.position.y < 0 {
		p.position.y = 0
	}

	if p.position.x > 1366 {
		p.position.x = 1366
	}

	/*
		if boots == true
		createpath() this is to create the ice walkway or whatever
	*/
}

func (p *Player) input() {
	p.state = PIdle

	//left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if p.velocity.x > -3 {
			p.velocity.x -= p.moveSpeed
		}
		p.state = PMovingLeft
	}

	//right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if p.velocity.x < 3 {
			p.velocity.x += p.moveSpeed
		}
		p.state = PMovingRight
	}

	//up
	if ebiten.IsKeyPressed(ebiten.KeyUp) && p.jump == false {
		if p.velocity.y < 1 {
			p.velocity.y -= p.moveSpeed * 50
		}
		p.state = PJumping
		p.jump = true
	}

	//"down"
	/*
		if p.state == PJumping {
			p.state = PFalling
			//rn treats 0 as the ground
			for p.position.y != 0 {
				p.position.y -= p.velocity.y
			}
			/*
				while onSolid false
					p.position.y -= p.velocity.y
	*/
	//}

	//movement abilities
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		abilities(p, "jax", "move")
		fmt.Println("skill used")
	}

	//stun test key
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.stunned = true
		go wait(p)

	}

	/*
		if press z{
			if you have frost runners {
				pressing z will dodge
				p.ability(p.type, dodge)
				set state to invincible
			}
			else{
				pressing z will have u use ur movement ability, p.ability(p.type, move)
				}
			}
		}

		if you press x && boots == false{
			if (ultRdy == true)
				p.ability(p.type, ult)
			else
				p.ability(p.type, atk)
				set state to hitting
		}

	*/

}

func wait(p *Player) {
	time.Sleep(3 * time.Second)
	p.stunned = false
}

func abilities(p *Player, character string, moveType string) {
	switch moveType {
	case "atk":

		switch character {

		case "jax":

		default:

		}

	case "move":

		switch character {

		case "jax":
			//extra jump
			if p.state == PMovingLeft {
				p.velocity.y = 0
				p.velocity.y -= 7
			}
		default:
			//dashes in direction moving
			if p.state == PMovingLeft {
				p.velocity.x = 0
				p.velocity.x -= 7
			} else if p.state == PMovingRight {
				p.velocity.x = 0
				p.velocity.x += 7
			}
		}

	//ults w/o boots
	case "ult":

		switch character {

		case "jax":

		default:

		}

	//dodging u need frostwalkers for this
	case "dodge":

		switch character {

		case "jax":

		default:

		}
	}
}

/*

abilities(character, type){

switch (type):

	case "atk":

		switch (character):

			case1: makes character play sprite animation and makes them hitting (will probably make hitbox sprite based)
			casecharacter2:
			casecharacter3:

	case "move":

		switch (character):

			case1: makes character play sprite animation and basically moves their position in whatever way we want
			case2:
			case3:

	case "ult":

		switch (character):

			case1: makes character ult however we want them to ult
			case2:
			case3:

	case "dodge":

		switch (character):

			case1: makes character dodge, giving i frames, their dodges can be different and act like move but w/ i frames
			case2:
			case3:
}

*/

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
