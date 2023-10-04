package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	Conn        net.Conn
	IsConnected bool
	ClientId    string
}

func RunTcpClient(port string) {
	data := make([]byte, 1024)
	host := fmt.Sprintf("localhost:%s", port)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	client := Client{
		Conn:        conn,
		IsConnected: true,
	}

	defer client.Conn.Close()
	client.Conn.Read(data)

	client.setClientId(data)

	go client.receiveMessage()
	client.readInputMessage()

	log.Println("connection was closed")
}
func (c *Client) readInputMessage() {
	in := bufio.NewReader(os.Stdin)

	for c.IsConnected {
		message, _ := in.ReadString('\n')
		payload := c.ClientId + "|" + message
		c.Conn.Write([]byte(payload))
	}
}
func (c *Client) receiveMessage() {
	for {
		data := make([]byte, 1024)
		_, err := c.Conn.Read(data)
		if err != nil {
			if errors.Is(err, net.ErrWriteToConnected) {
				continue
			}
			c.Conn.Close()
			c.IsConnected = false
			break
		}
		received := string(data)
		fmt.Println("server: ", received)
	}
}
func (c *Client) setClientId(idMsg []byte) {
	received := string(idMsg)
	clientId := strings.Split(received, ":")[1]
	clientId = strings.Split(clientId, "\n")[0]
	fmt.Println(received)
	c.ClientId = clientId
}
