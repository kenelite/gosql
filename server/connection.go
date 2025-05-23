package server

import (
	"github.com/kenelite/gosql/executor"
	"log"
	"net"
)

func (s *Server) handleConnection(conn net.Conn) {
	client := protocol.NewClientConn(conn)
	if err := client.Handshake(); err != nil {
		log.Println("Handshake failed:", err)
		return
	}

	for {
		query, err := client.ReadQuery()
		if err != nil {
			break
		}
		result := executor.Execute(query)
		client.WriteResult(result)
	}
}
