package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	data := make([]byte, 1024)
	in := bufio.NewReader(os.Stdin)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Read(data)

	received := string(data)
	clientId := strings.Split(received, ":")[1]
	clientId = strings.Split(clientId, "\n")[0]
	fmt.Println(received)

	var message string
	for {
		message, _ = in.ReadString('\n')
		payload := clientId + "|" + message
		conn.Write([]byte(payload))
		go receiveMessage(conn)
	}

}
func receiveMessage(conn net.Conn) {
	data := make([]byte, 1024)
	for {
		_, err := conn.Read(data)
		if err != nil {
			if errors.Is(err, net.ErrWriteToConnected) {
				continue
			}
			break
		}
		received := string(data)
		fmt.Println("server: ", received)
	}
}
