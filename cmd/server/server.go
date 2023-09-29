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
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	Id     string
	Comand string
	Body   string
}
type Server struct {
	TotalConns   int
	ConnectedIds []string
}

func main() {
	port := os.Getenv("PORT")
	idChannels := make(chan string)
	server := Server{}
	var wg sync.WaitGroup

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	log.Println("running on ", port)
	go server.ServerCounter(&wg, idChannels)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go server.HandleConnection(conn, idChannels)
	}
}
func (s *Server) ServerCounter(wg *sync.WaitGroup, ch chan string) {
	findIndex := func(slice []string, target string) int {
		for i, v := range slice {
			if v == target {
				return i
			}
		}
		return -1
	}
	wg.Add(1)
	defer wg.Done()
	for {
		newId := <-ch
		index := findIndex(s.ConnectedIds, newId)
		if index > -1 {
			s.TotalConns -= 1
			s.ConnectedIds = append(s.ConnectedIds[:index], s.ConnectedIds[index+1:]...)
			continue
		}
		s.TotalConns += 1
		s.ConnectedIds = append(s.ConnectedIds, newId)
		fmt.Println("Counter:", s.TotalConns)
	}
}
func (s *Server) HandleConnection(con net.Conn, ch chan string) {
	defer con.Close()
	id := s.generateID()
	defer func() {
		ch <- id
	}()
	ch <- id
	message := "welcome, your id: " + id + "\n"
	con.Write([]byte(message))
	reader := bufio.NewReader(con)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			if errors.Is(err, net.ErrWriteToConnected) {
				continue
			}
			break
		}
		msg, err := s.ParseMessage(message)

		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = s.CheckCommand(msg)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if msg.Comand == "LIST" {
			s.SendIds(con)
		}
	}
}

func (s *Server) ParseMessage(m string) (Message, error) {
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
func (s *Server) CheckCommand(m Message) error {
	if m.Comand != "LIST" && m.Comand != "RELAY" {
		return errors.New("invalid comand")
	}
	return nil
}
func (s *Server) SendIds(con net.Conn) {
	idMessage := strings.Join(s.ConnectedIds, ",")
	m := "Connected: " + idMessage
	con.Write([]byte(m))
}
func (s *Server) generateID() string {
	id := uuid.New().String()
	id = strings.Split(id, "-")[0]
	return id
}
