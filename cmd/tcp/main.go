package main

import (
	"fmt"
	"os"

	"github.com/LuisKpBeta/tcp-server/pkg/client"
	"github.com/LuisKpBeta/tcp-server/pkg/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT env must not be null")
	}
	args := os.Args
	if len(args) < 2 {
		fmt.Println("you must define 'server' or 'client' to run")
		return
	}
	execMode := args[1]
	if execMode != "server" && execMode != "client" {
		fmt.Println("you must define 'server' or 'client' to run")
		return
	}

	if execMode == "server" {
		server.CreateAndRunServer(port)
		return
	}
	if execMode == "client" {
		client.RunTcpClient(port)
		return
	}

}
