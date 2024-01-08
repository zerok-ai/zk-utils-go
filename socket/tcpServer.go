package socket

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"net"
)

const (
	LoggerTagServer = "TCP Socket Server"
)

type HandleTCPData func([]byte)

type TCPServer struct {
	HandleTCPData HandleTCPData
	Port          string
	listener      net.Listener
	connections   []net.Conn
}

func (server *TCPServer) handleConnection(conn net.Conn) {

	output := readData(conn)
	if server.HandleTCPData != nil {
		server.HandleTCPData(output)
	}

	// Echo the data back to the client
	//_, err := conn.Write(buffer[:n])
	//if err != nil {
	//	fmt.Println("Error writing:", err)
	//	return
	//}
}

func (server *TCPServer) CloseServer() {

	for _, conn := range server.connections {
		err := conn.Close()
		if err != nil {
			zkLogger.Error(LoggerTagServer, "Error closing connection:", err)
		}
	}

	err := server.listener.Close()
	if err != nil {
		zkLogger.Error(LoggerTagServer, "Error closing tcp listener:", err)
	}
}

func CreateTCPServer(port string, handleTCPData HandleTCPData) *TCPServer {
	return &TCPServer{Port: port, HandleTCPData: handleTCPData}
}

func (server *TCPServer) Start() {

	// Start listening
	listener, err := net.Listen("tcp", ":"+server.Port)
	if err != nil {
		zkLogger.Error(LoggerTagServer, "Error listening:", err)
		return
	}
	server.listener = listener
	server.connections = make([]net.Conn, 0)

	zkLogger.Info(LoggerTagServer, "Server is listening on port "+server.Port)

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
