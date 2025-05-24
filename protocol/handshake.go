package protocol

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"math/rand"
)

var serverVersion = "5.7.0-gosql"

func (c *Conn) Handshake(users map[string]string) error {
	// 1. Send handshake initialization
	authPluginData := make([]byte, 20)
	rand.Read(authPluginData)

	var buf bytes.Buffer
	buf.WriteByte(0x0a) // protocol version
	buf.WriteString(serverVersion)
	buf.WriteByte(0x00)                                     // null terminator
	binary.Write(&buf, binary.LittleEndian, uint32(1))      // connection ID
	buf.Write(authPluginData[:8])                           // auth-plugin-data-part-1
	buf.WriteByte(0x00)                                     // filler
	binary.Write(&buf, binary.LittleEndian, uint16(0x0002)) // capability flags (lower 2 bytes)
	buf.WriteByte(0x21)                                     // character set
	binary.Write(&buf, binary.LittleEndian, uint16(0))      // status flags
	binary.Write(&buf, binary.LittleEndian, uint16(0x8000)) // capability flags (upper 2 bytes)
	buf.WriteByte(20)                                       // length of auth-plugin-data
	buf.Write(make([]byte, 10))                             // reserved
	buf.Write(authPluginData[8:])                           // auth-plugin-data-part-2
	buf.WriteByte(0x00)                                     // null terminator for plugin name

	if err := c.WritePacket(buf.Bytes()); err != nil {
		return err
	}

	// 2. Read login request
	data, err := c.ReadPacket()
	if err != nil {
		return err
	}

	username, clientAuth := parseLogin(data, authPluginData)
	storedPass, ok := users[username]
	if !ok || !checkPassword(storedPass, authPluginData, clientAuth) {
		return c.WriteError(1045, "Access denied for user")
	}

	// 3. Send OK packet
	return c.WriteOK()
}

func parseLogin(data, seed []byte) (string, []byte) {
	pos := 36
	end := bytes.IndexByte(data[pos:], 0x00)
	if end == -1 {
		return "", nil
	}
	username := string(data[pos : pos+end])
	pos += end + 1

	if pos >= len(data) {
		return username, nil
	}
	authLen := int(data[pos])
	pos++
	if pos+authLen > len(data) {
		return username, nil
	}
	authResp := data[pos : pos+authLen]

	return username, authResp
}

func checkPassword(password string, seed, authResp []byte) bool {
	if password == "" {
		return len(authResp) == 0
	}

	shaPass := sha1.Sum([]byte(password))
	shaShaPass := sha1.Sum(shaPass[:])

	h := sha1.New()
	h.Write(seed[:20])
	h.Write(shaShaPass[:])
	hash := h.Sum(nil)

	if len(hash) != len(authResp) {
		return false
	}

	for i := range hash {
		hash[i] ^= authResp[i]
	}

	return bytes.Equal(hash, shaPass[:])
}
