package protocol

import (
	"bytes"
	"io"
)

type Conn struct {
	Reader io.Reader
	Writer io.Writer
	Seq    uint8
}

// ReadPacket reads a full MySQL packet (4-byte header + payload)
func (c *Conn) ReadPacket() ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(c.Reader, header); err != nil {
		return nil, err
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	c.Seq = header[3]

	data := make([]byte, length)
	if _, err := io.ReadFull(c.Reader, data); err != nil {
		return nil, err
	}
	return data, nil
}

// WritePacket writes a MySQL packet
func (c *Conn) WritePacket(data []byte) error {
	length := len(data)
	header := []byte{
		byte(length),
		byte(length >> 8),
		byte(length >> 16),
		c.Seq,
	}
	c.Seq++

	packet := append(header, data...)
	_, err := c.Writer.Write(packet)
	return err
}

// WriteEOF sends a MySQL EOF packet (used after column definitions and after rows)
func (c *Conn) WriteEOF() error {
	var buf bytes.Buffer
	buf.WriteByte(0xfe) // EOF packet header
	buf.WriteByte(0x00) // warnings (lower byte)
	buf.WriteByte(0x00) // warnings (upper byte)
	buf.WriteByte(0x00) // status flags (lower byte)
	buf.WriteByte(0x00) // status flags (upper byte)
	return c.WritePacket(buf.Bytes())
}
