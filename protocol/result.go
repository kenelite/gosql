package protocol

import (
	"bytes"
	"fmt"
	"github.com/kenelite/gosql/storage"
)

func (c *Conn) WriteResultSet(columns []string, rows []storage.Row) error {
	var err error

	// Send column count as a LengthEncodedInteger
	c.Seq++
	err = c.WritePacket(appendLengthEncodedInt(nil, uint64(len(columns))))
	if err != nil {
		return err
	}

	// Send column definitions
	for _, name := range columns {
		c.Seq++
		var buf bytes.Buffer
		buf.WriteByte(0x03)                       // Catalog: "def"
		writeLenEncString(&buf, "def")            // Catalog
		writeLenEncString(&buf, "")               // Schema
		writeLenEncString(&buf, "")               // Table
		writeLenEncString(&buf, "")               // OrgTable
		writeLenEncString(&buf, name)             // Name
		writeLenEncString(&buf, name)             // OrgName
		buf.WriteByte(0x0c)                       // Fixed length of next fields
		buf.Write([]byte{0x3f, 0x00})             // Charset (utf8_general_ci)
		buf.Write([]byte{0xff, 0xff, 0xff, 0xff}) // Column length
		buf.WriteByte(0x03)                       // Type: MYSQL_TYPE_LONG (int)
		buf.Write([]byte{0, 0})                   // Flags
		buf.WriteByte(0)                          // Decimals
		buf.Write([]byte{0, 0})                   // Filler
		err = c.WritePacket(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// EOF packet
	c.Seq++
	err = c.WriteEOF()
	if err != nil {
		return err
	}

	// Send rows
	for _, row := range rows {
		c.Seq++
		var buf bytes.Buffer
		for _, val := range row {
			switch v := val.(type) {
			case nil:
				buf.WriteByte(0xfb) // NULL
			case int:
				writeLenEncString(&buf, fmt.Sprint(v))
			case string:
				writeLenEncString(&buf, v)
			default:
				writeLenEncString(&buf, fmt.Sprintf("%v", v))
			}
		}
		err = c.WritePacket(buf.Bytes())
		if err != nil {
			return err
		}
	}

	// Final EOF
	c.Seq++
	return c.WriteEOF()
}
