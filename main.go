package main

import (
	"github.com/kenelite/gosql/protocol"
	"github.com/kenelite/gosql/storage"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":3306")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	store, err := storage.NewFileStore()
	if err != nil {
		log.Fatal("failed to initialize storage:", err)
	}

	log.Println("gosql listening on :3306")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConnection(conn, store)
	}
}

func handleConnection(nc net.Conn, store *storage.FileStore) {
	defer nc.Close()
	c := &protocol.Conn{Reader: nc, Writer: nc}

	// Initial handshake
	if err := c.WriteHandshake(); err != nil {
		log.Println("handshake error:", err)
		return
	}

	_, err := c.ReadPacket() // Login request
	if err != nil {
		log.Println("login read error:", err)
		return
	}
	c.Seq = 2
	c.WriteOK()

	// Handle client queries
	for {
		data, err := c.ReadPacket()
		if err != nil {
			return // client disconnected
		}

		if len(data) == 0 || data[0] != 0x03 {
			c.Seq = 1
			c.WriteError(1064, "Only COM_QUERY supported")
			continue
		}

		query := string(data[1:])
		log.Println("Query:", query)
		c.Seq = 1
		handleQuery(c, store, query)
	}
}

func handleQuery(c *protocol.Conn, store *storage.FileStore, query string) {
	query = strings.TrimSpace(query)
	queryLower := strings.ToLower(query)

	switch {
	case queryLower == "select 1":
		c.WriteResultSet([]string{"1"}, [][]interface{}{{1}})

	case strings.HasPrefix(queryLower, "create table"):
		// Minimal parser: CREATE TABLE users (id INT, name VARCHAR)
		parts := strings.Split(query, "(")
		if len(parts) != 2 {
			c.WriteError(1064, "invalid CREATE TABLE syntax")
			return
		}
		name := strings.Fields(parts[0])[2]
		colsDef := strings.TrimSuffix(parts[1], ")")
		colParts := strings.Split(colsDef, ",")
		cols := []storage.Column{}
		for _, col := range colParts {
			col = strings.TrimSpace(col)
			tokens := strings.Fields(col)
			if len(tokens) < 2 {
				c.WriteError(1064, "invalid column definition")
				return
			}
			cols = append(cols, storage.Column{Name: tokens[0], Type: tokens[1]})
		}
		if err := store.CreateTable(name, cols); err != nil {
			c.WriteError(1064, err.Error())
			return
		}
		c.WriteOK()

	case strings.HasPrefix(queryLower, "insert into"):
		// Minimal parser: INSERT INTO users VALUES (1, 'Alice')
		parts := strings.Split(query, "values")
		if len(parts) != 2 {
			c.WriteError(1064, "invalid INSERT syntax")
			return
		}
		table := strings.Fields(parts[0])[2]
		valsRaw := strings.TrimSpace(parts[1])
		valsRaw = strings.Trim(valsRaw, "()")
		valTokens := strings.Split(valsRaw, ",")
		row := make([]interface{}, 0, len(valTokens))
		for _, tok := range valTokens {
			tok = strings.TrimSpace(tok)
			if strings.HasPrefix(tok, "'") && strings.HasSuffix(tok, "'") {
				row = append(row, strings.Trim(tok, "'"))
			} else {
				row = append(row, parseInt(tok))
			}
		}
		if err := store.Insert(table, row); err != nil {
			c.WriteError(1064, err.Error())
			return
		}
		c.WriteOK()

	case strings.HasPrefix(queryLower, "select"):
		// Minimal parser: SELECT * FROM users
		tokens := strings.Fields(queryLower)
		if len(tokens) < 4 || tokens[1] != "*" || tokens[2] != "from" {
			c.WriteError(1064, "unsupported SELECT syntax")
			return
		}
		table := tokens[3]
		cols, rows, err := store.SelectAll(table)
		if err != nil {
			c.WriteError(1064, err.Error())
			return
		}
		names := make([]string, len(cols))
		for i, col := range cols {
			names[i] = col.Name
		}
		c.WriteResultSet(names, rows)

	default:
		c.WriteError(1064, "unsupported query")
	}
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
