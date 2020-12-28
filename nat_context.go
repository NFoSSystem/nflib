package nflib

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

func ipv4ToInt32(ip *net.IP) uint32 {
	ipB := []byte(*ip)
	return uint32(ipB[12])<<24 | uint32(ipB[13])<<16 | uint32(ipB[14])<<8 | uint32(ipB[15])
}

func uint32ToIpv4(addr uint32) *net.IP {
	ip := net.IPv4(byte(addr>>24), byte((0xFFFFFF&addr)>>16), byte((0xFFFF&addr)>>8), byte(255&addr))
	return &ip
}

type Packet struct {
	Addr  uint32
	Crc16 uint16
	Port  uint16
}

func GetPacketFromBytes(bSlice []byte) *Packet {
	buff := &bytes.Buffer{}
	res := Packet{}
	buff.Grow(len(bSlice))
	_, err := buff.Write(bSlice)
	if err != nil {
		panic(err)
	}
	err = binary.Read(buff, binary.BigEndian, &res)
	if err != nil {
		panic(err)
	}
	return &res
}

func GetBytesFromPacket(src Packet) []byte {
	buff := &bytes.Buffer{}
	err := binary.Write(buff, binary.BigEndian, src)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func SendToRouter(pktChan chan Packet, addr *net.IP, port uint16) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{*addr, int(port), ""})
	if err != nil {
		log.Fatalf("Error opening TCP socket to %s:%d: %s\n", addr.String(), port)
	}
	defer conn.Close()

	for {
		select {
		case pkt := <-pktChan:
			_, err := conn.Write(GetBytesFromPacket(pkt))
			if err != nil {
				log.Fatalf("Error writing mapping message to %s:%d: %s", addr.String(), port, err)
			}
		default:
			continue
		}
	}
}
