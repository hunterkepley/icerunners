package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/hajimehoshi/ebiten"

	Packet "github.com/hunterkepley/defterra/packet"
)

// GameEntities holds all of the game's objects
type GameEntities struct {
	resources  GameResources // Resources
	gameMap    Map           // Game map
	camera     Camera        // Camera
	player     Player        // The actual player
	conn       *net.UDPConn
	buf        []byte
	chN        chan int
	chAddr     chan net.Addr
	serverAddr net.Addr
}

// GameResources holds all of the game's resources (images/music/sound)
type GameResources struct {
	//platformImages     PlatformImages
}

func (g *GameResources) init() {
	// Init resources
	//g.platformImages = createPlatformImages()
}

func (g *GameEntities) init() {
	// Initialize game resources
	g.resources.init()

	// Init map
	g.gameMap = initializeMap()
	g.camera = createCamera(Vec2f{0, 0}, 4)

	// Server (temporary until server browser/system made)
	g.serverAddr = &net.UDPAddr{IP: net.IPv4(18, 208, 230, 70), Port: 10001, Zone: ""}
	var err error
	g.conn, err = net.ListenUDP(
		"udp",
		&net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		fmt.Println(err)
	}

	// Packet info: 00000000 00001111
	ds := make([]byte, 2)
	binary.BigEndian.PutUint16(ds, 15)
	packet := Packet.CreatePacket(ds, []byte{0, 0}, []byte{}) // 0000000000001111 -- PCServerJoined (15)
	if err != nil {
		fmt.Println(err)
	}

	packet.Send(g.conn, g.serverAddr)

}

func updateGame(g *GameEntities) {

	// Listen to server
	listenToServer(g)

	// Update game map
	g.gameMap.update(g.camera)

	// Update entities
	g.player.update()

	// Update camera
	g.camera.update()
	g.gameMap.cameraPosition = g.camera.position
}

func drawGame(g *GameEntities, screen *ebiten.Image) {

	g.gameMap.clearImage()

	// Render map entities ---------------------------------
	g.player.render(g.gameMap.image)

	// Game map ---------------------------------
	g.gameMap.render(screen)

}

func listenToServer(g *GameEntities) {
	go func() {
		n, addr, err := g.conn.ReadFrom(g.buf) // Establish connection to client
		if err != nil {
			fmt.Println(err)
		}
		g.chN <- n
		g.chAddr <- addr
	}()

	select {
	case _ = <-g.chN:
		addr := <-g.chAddr
		g.serverAddr = addr
	case <-time.After(time.Microsecond):
	}

	p := Packet.DecodePacket(g.buf)

	if len(p.Code) > 0 {
		switch {
		case p.Code[0] == byte(0), p.Code[1] == byte(15):
			g.conn.Close()
		case p.Code[0] == Packet.DCameraZoom[0], p.Code[1] == Packet.DCameraZoom[1]:
			fmt.Println("e")
			g.camera.zoom = Packet.Byte2Float64(p.Data)
		}
	}
}
