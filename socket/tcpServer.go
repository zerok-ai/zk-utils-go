package socket

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"net"
)

type HandleTCPData func([]byte) string

type TCPServer struct {
	HandleTCPData HandleTCPData
	Port          string
	listener      *net.Listener
	connections   []net.Conn
	SendAck       bool
}

func (server *TCPServer) handleConnection(conn net.Conn) {

	for {
		output := readData(conn)
		fmt.Printf("Received on server: %s\n", output)
		var status string
		if server.HandleTCPData != nil {
			status = server.HandleTCPData(output)
		}

		//Echo the data back to the client
		if server.SendAck {
			_, err := conn.Write([]byte(status))
			if err != nil {
				fmt.Println("Error writing:", err)
				return
			}
		}
	}
}

func (server *TCPServer) Close() {

	if server.listener == nil {
		return
	}

	if server.connections == nil {
		return
	}

	for _, conn := range server.connections {
		err := conn.Close()
		if err != nil {
			zkLogger.Error(LoggerTagSocket, "Error closing connection:", err)
		}
	}

	err := (*server.listener).Close()
	if err != nil {
		zkLogger.Error(LoggerTagSocket, "Error closing tcp listener:", err)
	}
}

func CreateTCPServer(port string, handleTCPData HandleTCPData, sendAck bool) *TCPServer {
	return &TCPServer{Port: port, HandleTCPData: handleTCPData, SendAck: sendAck}
}

func (server *TCPServer) Start() {

	// Start listening
	listener, err := net.Listen("tcp", ":"+server.Port)
	if err != nil {
		zkLogger.Error(LoggerTagSocket, "Error listening:", err)
		return
	}
	server.listener = &listener
	server.connections = make([]net.Conn, 0)

	zkLogger.Info(LoggerTagSocket, "Server is listening on port "+server.Port)

	for {
		// Accept a connection
		conn, err1 := listener.Accept()
		if err1 != nil {
			fmt.Println("Error accepting connection:", err1)
			continue
		}
		server.connections = append(server.connections, conn)

		// Handle the connection in a new goroutine
		go server.handleConnection(conn)
	}
}
