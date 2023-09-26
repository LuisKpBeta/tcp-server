package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	log.Println("running on ", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn)
	}
}

func generateID() string {
	id := uuid.New().String()
	id = strings.Split(id, "-")[0]
	return id
}
func handleConnection(con net.Conn) {
	id := generateID()
	message := "welcome, you're " + id
	con.Write([]byte(message))
	con.Close()
}
