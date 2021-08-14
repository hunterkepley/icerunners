package Packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
)

/*
 *                           P A C K E T  S T R U C T U R E
 *
 *      Each packet has 4 bit code, then everything following is 4 bit datatype and data
 *                           corresponding to the data type.
 *
 *                    |      2 bytes      |    2 bytes     |  X bytes |
 *                       packet code int     datatype int      data
 *
 *
 */

// Packet sizes
const (
	PacketCodeSize     = 2 // bytes
	PacketDataTypeSize = 2 // bytes
)

// Data types
var (
	DCameraZoom     = []byte{0, 3} // 0000 0000 0000 0011
	DPlayerPosition = []byte{0, 4} // 0000 0000 0000 0100
)

type Packet struct {
	Code     []byte // Code for packet (what is is!)
	DataType []byte // All data types, parallel to v
	Data     []byte // All data in packet, parallel to ^
}

func CreatePacket(code []byte, dataType []byte, data []byte) Packet {
	return Packet{
		Code:     code,
		DataType: dataType,
		Data:     data,
	}
}

func (p *Packet) Send(ServerConn *net.UDPConn, addr net.Addr) {
	bs := make([]byte, 0) // Make packet byte slice -- code + datatype + X (data)

	bs = append(bs, p.Code...)     // Add Code
	bs = append(bs, p.DataType...) // Add Datatype
	bs = append(bs, p.Data...)     // Add Data

	_, err := ServerConn.WriteTo(bs, addr)
	if err != nil {
		fmt.Println("DS >> Packet failed to send to client at [", addr.String(), "]\n", err)
	}
}

func DecodePacket(p []byte) Packet {
	var codeBytes, dataTypeBytes, data []byte

Decode:
	for i, b := range p {
		switch i {
		case 0:
			codeBytes = append(codeBytes, b)
		case 1:
			codeBytes = append(codeBytes, b)
		case 2:
			dataTypeBytes = append(dataTypeBytes, b)
		case 3:
			dataTypeBytes = append(dataTypeBytes, b)
		default:
			data = p[i:] // Add data bytes
			break Decode // Exit loop
		}
	}

	if len(codeBytes) > 0 {
		// Check if there is a proper data type (0 is nonexistant)
		if codeBytes[0] == 0 && codeBytes[1] == 0 {
			return Packet{}
		}
	}

	return Packet{
		Code:     codeBytes,
		DataType: dataTypeBytes,
		Data:     data,
	}
}

func (p *Packet) String() string {
	return fmt.Sprintf("%d%d%b", p.Code, p.DataType, p.Data)
}

var ErrRange = errors.New("value out of range")

func BitStringToBytes(s string) []byte {
	var out []byte
	var str string

	for i := len(s); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(s[0:i])
		} else {
			str = string(s[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}

// Converts bytes into a 64-bit float
func Byte2Float64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// Converts Floats into a 8 bytes
func Float642Byte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func Float32fromBytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
