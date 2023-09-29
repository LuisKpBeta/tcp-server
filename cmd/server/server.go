package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type Message struct {
	Id     string
	Comand string
	Body   string
}

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
	defer con.Close()
	id := generateID()
	message := "welcome, your id: " + id + "\n"
	con.Write([]byte(message))
	reader := bufio.NewReader(con)
	// process
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			if errors.Is(err, net.ErrWriteToConnected) {
				continue
			}
			break
		}
		msg, err := HandleMessage(message)

		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = checkMessage(msg)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		m := fmt.Sprint("voce mandou:", msg.Comand)
		con.Write([]byte(m))

		// con.Write([]byte(""))
	}
	//
}

func HandleMessage(m string) (Message, error) {
	m = strings.Trim(m, "\n")
	m = strings.Trim(m, " ")
	pattern := `^([a-zA-Z0-9_]+)\|([a-zA-Z0-9_]+)(?::(.+))?$`
	r := regexp.MustCompile(pattern)
	matches := r.FindStringSubmatch(m)
	if len(matches) < 3 {
		return Message{}, errors.New("invalid message")
	}
	msg := Message{
		Id:     matches[1],
		Comand: matches[2],
	}
	if len(matches) == 3 {
		msg.Body = matches[3]
	}
	return msg, nil
}
func checkMessage(m Message) error {
	if m.Comand != "LIST" && m.Comand != "RELAY" {
		return errors.New("invalid comand")
	}
	return nil
}
