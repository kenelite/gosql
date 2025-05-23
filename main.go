package main

import (
	"log"
	"net"

	"github.com/kenelite/gosql/protocol"
)

func main() {
	listener, err := net.Listen("tcp", ":3306")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, _ := listener.Accept()
		go handle(conn)
	}
}

func handle(nc net.Conn) {
	defer nc.Close()
	c := &protocol.Conn{
		Reader: nc,
		Writer: nc,
	}

	if err := c.WriteHandshake(); err != nil {
		log.Println("handshake err:", err)
		return
	}

	data, err := c.ReadPacket() // Login Request
	if err != nil {
		return
	}

	log.Println("Login packet:", data)

	c.Seq = 2
	c.WriteOK()

	for {
		data, err := c.ReadPacket()
		if err != nil {
			return
		}
		log.Printf("Query packet: %q\n", data)
		c.Seq = 1
		c.WriteError(1064, "Only SELECT 1 supported")
	}
}
