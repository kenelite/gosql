package protocol

import (
	"bytes"
	"crypto/rand"
)

const (
	protocolVersion = 10
	serverVersion   = "5.7.0-gosql"
)

func (c *Conn) WriteHandshake() error {
	buf := &bytes.Buffer{}
	buf.WriteByte(protocolVersion)
	buf.WriteString(serverVersion)
	buf.WriteByte(0x00) // Null-terminated

	// Connection ID
	buf.Write([]byte{1, 0, 0, 0})

	// Auth plugin data part 1
	salt := make([]byte, 8)
	rand.Read(salt)
	buf.Write(salt)
	buf.WriteByte(0x00)

	// Capability flags (lower 2 bytes)
	buf.Write([]byte{0xff, 0xf7}) // CLIENT_PROTOCOL_41, CLIENT_SECURE_CONN
	buf.WriteByte(0x21)           // Charset: utf8
	buf.Write([]byte{0x00, 0x00}) // Status flags
	buf.Write([]byte{0xff, 0xff}) // Capability flags (upper 2 bytes)
	buf.WriteByte(21)             // Auth plugin length
	buf.Write(make([]byte, 10))   // Reserved

	// Auth plugin data part 2
	salt2 := make([]byte, 12)
	rand.Read(salt2)
	buf.Write(salt2)
	buf.WriteByte(0x00)

	// Auth plugin name
	buf.WriteString("mysql_native_password")
	buf.WriteByte(0x00)

	return c.WritePacket(buf.Bytes())
}
