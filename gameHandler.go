package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"

	Packet "github.com/hunterkepley/defterra/packet"
)

// GameEntities holds all of the game's objects
type GameEntities struct {
	resources GameResources // Resources
	gameMap   Map           // Game map
	camera    Camera        // Camera
	player    Player        // The actual player

	// Networking stuff (keep below other entities)
	conn               *net.UDPConn
	buf                []byte
	chN                chan int
	chAddr             chan net.Addr
	serverAddr         net.Addr
	wgSendDataToServer sync.WaitGroup
	id                 []byte
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
	g.serverAddr = &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 10001, Zone: ""}
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

	// Send player data to server
	go sendPlayerDataToServer(g, 32, Packet.DPlayerPositionX[:], g.player.position.x)
	g.wgSendDataToServer.Add(1)
	go sendPlayerDataToServer(g, 32, Packet.DPlayerPositionY[:], g.player.position.y)
	g.wgSendDataToServer.Add(1)

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

func sendPlayerDataToServer(g *GameEntities, _code uint16, _type []byte, _data float64) {
	defer g.wgSendDataToServer.Done()
	// Packet info: 0000 0000 code
	// Player position
	ds := make([]byte, 2)                 // 0000 0000
	binary.BigEndian.PutUint16(ds, _code) // 0000 0000 + code0:3 code4:7

	data := make([]byte, 16)
	binary.LittleEndian.PutUint64(data[:], math.Float64bits(_data))
	packet := Packet.CreatePacket(ds, _type, data[:])

	packet.Send(g.conn, g.serverAddr)
	// ^ send player position
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
		case p.Code[0] == Packet.DPlayerPositionX[0], p.Code[1] == Packet.DPlayerPositionX[1]:
		}
	}
}
