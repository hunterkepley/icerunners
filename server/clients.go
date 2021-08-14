package main

import (
	"net"

	Packet "github.com/hunterkepley/defterra/packet"
)

type Client struct {
	addr                        net.Addr
	connected                   bool
	timeSinceLastPacketReceived int
	timeSinceLastPacketSent     int
}

func sendPacketToAllClients(packet *Packet.Packet, ServerConn *net.UDPConn, clients []Client) {
	for i := 0; i < len(clients); i++ {
		if i >= len(clients) {
			break
		}

		packet.Send(ServerConn, clients[i].addr)
		clients[i].timeSinceLastPacketReceived = 0
	}
}
