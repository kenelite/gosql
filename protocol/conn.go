package protocol

import (
	"bufio"
	"errors"
	"fmt"
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
	data, err := c.ReadPacket()
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New("empty packet")
	}

	switch data[0] {
	case 0x01: // COM_QUIT
		return "", ErrQuit
	case 0x03: // COM_QUERY
		return string(data[1:]), nil
	default:
		return "", fmt.Errorf("unsupported command: 0x%x", data[0])
	}
}

func (c *Conn) ReadPacket() ([]byte, error) {
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
