# TCP Server


# Requirements

- [x] Create a TCP server that listens on a port specified by an environment variable PORT.
- [x] When a client connects to the server, a unique ID is generated and sent back to the client
- [x] Quando um cliente se conecta ao servidor, um id único é gerado e enviado de volta ao cliente
- [x] The client can send messages with three parts of data:
  - id - the client's ID (mandatory)
  - action - an action expected by the server (mandatory)
  - body - any data (JSON, plain text, anything) that a client wants to send to the server (optional).

- [x] When the client sends a message with the action LIST, the server should return a list of all connected client IDs
- [x] When the client sends a message with the action RELAY, the server should send the message from the body field to all connected clients.


# How to use
## Build
For the TCP client/server project, we have uploaded a vendor folder with dependencies. Therefore, you will need to build the package using:

 `go build cmd/tcp/main.go`

## Run
After building the Go package, you can use this package by running one of the following commands:
- `PORT=8080 ./main server` to run the server on PORT 8080.
- `PORT=8080 ./main client`  to run the client, connecting to the server on localhost, which is running on port 8080.

## Commands
The available commands are:
- `LIST` to get ID for all connected clients
- `RELAY:<message>` to send `<message>` to all clients

If you are not connected using the official client, you must send a message using this pattern:

`<id>|<comand>:<optional_message>`