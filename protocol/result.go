package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/kenelite/gosql/storage"
)

func (c *Conn) WriteResultSet(columns []string, rows []storage.Row) error {
	var err error

	// Column count
	err = c.WritePacket(encodeLenEncInt(uint64(len(columns))))
	if err != nil {
		return err
	}

	// Column definitions
	for _, col := range columns {
		err = c.WritePacket(encodeColumnDef(col))
		if err != nil {
			return err
		}
	}

	// EOF after columns
	err = c.WriteEOF()
	if err != nil {
		return err
	}

	// Rows
	for _, row := range rows {
		buf := new(bytes.Buffer)
		for _, val := range row {
			var s string
			if val == nil {
				s = "NULL"
			} else {
				s = fmt.Sprintf("%v", val) // convert val (interface{}) to string
			}
			buf.Write(encodeLenEncString(s))
		}
		err = c.WritePacket(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// EOF after rows
	return c.WriteEOF()
}

func (c *Conn) WriteEOF() error {
	packet := []byte{0xfe, 0x00, 0x00, 0x00, 0x00} // EOF header + status flags + warnings
	return c.WritePacket(packet)
}

func encodeLenEncInt(n uint64) []byte {
	switch {
	case n < 251:
		return []byte{byte(n)}
	case n < 1<<16:
		return []byte{0xfc, byte(n), byte(n >> 8)}
	case n < 1<<24:
		return []byte{0xfd, byte(n), byte(n >> 8), byte(n >> 16)}
	default:
		buf := make([]byte, 9)
		buf[0] = 0xfe
		binary.LittleEndian.PutUint64(buf[1:], n)
		return buf
	}
}

func encodeLenEncString(s string) []byte {
	length := encodeLenEncInt(uint64(len(s)))
	return append(length, s...)
}

func encodeColumnDef(name string) []byte {
	// minimal implementation
	buf := new(bytes.Buffer)
	buf.Write(encodeLenEncString("def"))                   // catalog
	buf.Write(encodeLenEncString("gosql"))                 // schema
	buf.Write(encodeLenEncString(""))                      // table
	buf.Write(encodeLenEncString(""))                      // org_table
	buf.Write(encodeLenEncString(name))                    // name
	buf.Write(encodeLenEncString(name))                    // org_name
	buf.WriteByte(0x0c)                                    // fixed length fields
	binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // charset
	binary.Write(buf, binary.LittleEndian, uint32(256))    // column length
	buf.WriteByte(0xfd)                                    // type = VARCHAR
	buf.WriteByte(0x00)                                    // flags
	buf.WriteByte(0x00)                                    // decimals
	buf.WriteByte(0x00)                                    // filler
	buf.WriteByte(0x00)                                    // filler
	return buf.Bytes()
}
