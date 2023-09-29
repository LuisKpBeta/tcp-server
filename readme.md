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


 