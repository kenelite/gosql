package protocol

import "bytes"

func (c *Conn) WriteOK() error {
	buf := &bytes.Buffer{}
	buf.WriteByte(0x00)           // OK packet header
	buf.WriteByte(0x00)           // affected rows
	buf.WriteByte(0x00)           // last insert ID
	buf.Write([]byte{0x00, 0x00}) // status flags
	buf.Write([]byte{0x00, 0x00}) // warnings
	return c.WritePacket(buf.Bytes())
}

func (c *Conn) WriteError(code int, msg string) error {
	buf := &bytes.Buffer{}
	buf.WriteByte(0xff)
	buf.Write([]byte{byte(code), byte(code >> 8)})
	buf.WriteByte('#')
	buf.WriteString("HY000")
	buf.WriteString(msg)
	return c.WritePacket(buf.Bytes())
}
