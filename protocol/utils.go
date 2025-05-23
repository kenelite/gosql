package protocol

import (
	"bytes"
	"encoding/binary"
)

func appendLengthEncodedInt(dst []byte, n uint64) []byte {
	switch {
	case n < 251:
		return append(dst, byte(n))
	case n < 1<<16:
		return append(dst, 0xfc, byte(n), byte(n>>8))
	case n < 1<<24:
		return append(dst, 0xfd, byte(n), byte(n>>8), byte(n>>16))
	default:
		b := make([]byte, 9)
		b[0] = 0xfe
		binary.LittleEndian.PutUint64(b[1:], n)
		return append(dst, b...)
	}
}

func writeLenEncString(buf *bytes.Buffer, s string) {
	buf.Write(appendLengthEncodedInt(nil, uint64(len(s))))
	buf.WriteString(s)
}
