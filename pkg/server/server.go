package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
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
type Connection struct {
	Con    net.Conn
	Id     string
	Active bool
}
type Server struct {
	TotalConns  int
	Connections []*Connection
}

func CreateAndRunServer(port string) {
	server := Server{}
	idChannels := make(chan *Connection)
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

func (s *Server) ServerCounter(wg *sync.WaitGroup, ch chan *Connection) {
	findIndex := func(slice []*Connection, target string) int {
		for i, v := range slice {
			if v.Id == target {
				return i
			}
		}
		return -1
	}
	wg.Add(1)
	defer wg.Done()
	for {
		con := <-ch
		index := findIndex(s.Connections, con.Id)
		if index > -1 {
			s.TotalConns -= 1
			s.Connections = append(s.Connections[:index], s.Connections[index+1:]...)
			continue
		}
		s.TotalConns += 1
		s.Connections = append(s.Connections, con)
	}
}
func (s *Server) HandleConnection(con net.Conn, ch chan *Connection) {
	defer con.Close()
	id := s.generateID()
	newCon := Connection{
		Id:     id,
		Active: true,
		Con:    con,
	}
	defer func() {
		newCon.Active = false
		ch <- &newCon
	}()

	ch <- &newCon
	message := "welcome, your id: " + id + "\n"
	log.Printf("client %s connected", id)
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
		log.Printf("client %s send %s", id, msg.Comand)
		if msg.Comand == "LIST" {
			s.SendIds(con)
		}
		if msg.Comand == "RELAY" {
			s.SendMessageForAll(msg.Body)
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
	if len(matches) == 4 {
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
	var connectionIDs string

	for i, conn := range s.Connections {
		connectionIDs += conn.Id
		if i < len(s.Connections)-1 {
			connectionIDs += ", "
		}
	}
	m := "Connected: " + connectionIDs
	con.Write([]byte(m))
}
func (s *Server) SendMessageForAll(msg string) {
	for _, conn := range s.Connections {
		if conn.Active && conn.Con != nil {

			conn.Con.Write([]byte(msg))
		}
	}
}
func (s *Server) generateID() string {
	id := uuid.New().String()
	id = strings.Split(id, "-")[0]
	return id
}
