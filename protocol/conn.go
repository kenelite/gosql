package protocol

import (
	"bufio"
	"io"
	"net"
)

type Conn struct {
	netConn net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	seq     uint8
}

func NewConn(nc net.Conn) *Conn {
	return &Conn{
		netConn: nc,
		reader:  bufio.NewReader(nc),
		writer:  bufio.NewWriter(nc),
		seq:     0,
	}
}

func (c *Conn) ReadQuery() (string, error) {
	data, err := c.readPacket()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Conn) readPacket() ([]byte, error) {
	head := make([]byte, 4)
	if _, err := io.ReadFull(c.reader, head); err != nil {
		return nil, err
	}
	length := int(head[0]) | int(head[1])<<8 | int(head[2])<<16
	c.seq = head[3]
	body := make([]byte, length)
	if _, err := io.ReadFull(c.reader, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Conn) WritePacket(data []byte) error {
	length := len(data)
	head := []byte{byte(length), byte(length >> 8), byte(length >> 16), c.seq}
	c.seq++
	packet := append(head, data...)
	_, err := c.writer.Write(packet)
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Conn) Close() {
	c.netConn.Close()
}
