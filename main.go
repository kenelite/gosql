package main

import (
	"fmt"
	"github.com/kenelite/gosql/config"
	"github.com/kenelite/gosql/executor"
	"github.com/kenelite/gosql/parser"
	"github.com/kenelite/gosql/protocol"
	"github.com/kenelite/gosql/storage"
	"net"
	"os"
)

func main() {
	cfg := config.Load()

	store, err := storage.NewFileStore(cfg.DataDir)
	if err != nil {
		fmt.Println("Failed to initialize storage:", err)
		os.Exit(1)
	}

	exec := executor.NewExecutor(store)

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		os.Exit(1)
	}
	fmt.Println("gosql listening on", cfg.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn, exec)
	}
}

func handleConnection(nc net.Conn, exec *executor.Executor) {
	conn := protocol.NewConn(nc)
	defer conn.Close()

	for {
		query, err := conn.ReadQuery()
		if err != nil {
			return
		}

		stmt, err := parser.Parse(query)
		if err != nil {
			conn.WriteError(1064, err.Error())
			continue
		}

		err = exec.Execute(stmt, conn)
		if err != nil {
			fmt.Println("Execution error:", err)
			return
		}
	}
}
