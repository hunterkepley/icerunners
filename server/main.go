package main

import (
	"fmt"
	"net"
	"time"

	Packet "github.com/hunterkepley/defterra/packet"
)

const (
	MaxDisconnectTime int = 1000 // How long no packets are sent until client is considered disconnected

	TickRate = 20 // Tickrate!
)

var (
	clients    []Client
	totalTicks int
)

// Packet codes:
var (
	// Server packet codes (prefixed by 0)
	PCServerJoined = []byte{0, 15} // 0000 0000 0000 1111
	PCServerData   = []byte{0, 1}  // 0000 0000 0000 0001 -- Just data and datatypes, no special actions
	// Client packet codes (prefixed by 1)
	PCClientJoined = []byte{0, 16} // 0000 0000 0001 0000
	PCClientData   = []byte{0, 32} // 0000 0000 0010 0000 -- Just data and datatypes, no special actions
)

func main() {
	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 10001})
	if err != nil {
		fmt.Println("DS >> ", err)
	}
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	chN := make(chan int)
	chAddr := make(chan net.Addr)

	for { // Server loop
		go func() {
			n, addr, err := ServerConn.ReadFrom(buf) // Establish connection to client
			chAddr <- addr
			chN <- n
			if err != nil {
				fmt.Println("DS >> Failed to read from UDP to buf\n--\t ", err)
			}
		}()

		var addr net.Addr
		var noConnection bool
		select {
		case addr = <-chAddr:
			// Received packet!
		case <-time.After(time.Second / TickRate):
			noConnection = true
		}
		if !noConnection {
			packet := Packet.DecodePacket(buf)
			if len(packet.Code) == 2 {
				switch { // Check if basic packet
				case packet.Code[0] == PCServerJoined[0] && packet.Code[1] == PCServerJoined[1]:
					contains := false
					for i := 0; i < len(clients); i++ { // Check if client exists
						if clients[i].connected && clients[i].addr.String() == addr.String() {
							clients[i].timeSinceLastPacketSent = 0
							contains = true
							// TODO: Tell client that they are already connected to the server so they cannot connect on another client
						}
					}
					if !contains {
						fmt.Println(addr)
						clients = append(clients, Client{addr, true, 0, 0}) // If doesn't exist, add to server list!
						fmt.Println("DS >>", len(clients), " clients connected [", addr.String(), "]")
					}
				}
			}
		}

		// Send test data
		d := Packet.Float642Byte(0.5)
		p := Packet.CreatePacket(PCClientData, Packet.DCameraZoom, d)
		sendPacketToAllClients(&p, ServerConn, clients)

		for i := 0; i < len(clients); i++ {
			clients[i].timeSinceLastPacketSent++
			if noConnection {
				clients[i].timeSinceLastPacketSent += 50
			}
			clients[i].timeSinceLastPacketReceived++

			// Disconnect client for timeout
			if clients[i].timeSinceLastPacketSent > MaxDisconnectTime {
				clients[i].connected = false
				fmt.Println("DS >> Client at [", clients[i].addr.String(), "] Disconnected [Timed out...]")
				clients[i] = clients[len(clients)-1] // Delete dc'd client from slice
				clients[len(clients)-1] = Client{}
				clients = clients[:len(clients)-1] // Delete dc'd client from slice
				break
			}
		}

		totalTicks++
		// Print for debugging
		if totalTicks%1000 == 0 {
			fmt.Println("DS >> Total ticks:", totalTicks)
		}
	}
}
