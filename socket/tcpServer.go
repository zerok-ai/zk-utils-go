package socket

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"net"
)

const (
	LoggerTag = "TCP Socket"
)

type HandleTCPData func([]byte)

type TCPServer struct {
	HandleTCPData HandleTCPData
	Port          string
}

func (server TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	output := make([]byte, 0)

	for {
		// Read data from the connection
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		if n == 0 {
			break
		}

		// Print the received data
		fmt.Printf("Received: %s", buffer[:n])

		output = append(output, buffer[:n]...)

	}

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

func CreateTCPServer(port string, handleTCPData HandleTCPData) *TCPServer {
	return &TCPServer{Port: port, HandleTCPData: handleTCPData}
}

func (server TCPServer) Start() {

	// Start listening
	listener, err := net.Listen("tcp", ":"+server.Port)
	if err != nil {
		zkLogger.Error(LoggerTag, "Error listening:", err)
		return
	}

	defer func(listener net.Listener) {
		err = listener.Close()
		if err != nil {
			zkLogger.Error(LoggerTag, "Error closing tcp listener:", err)
		}
	}(listener)

	zkLogger.Info(LoggerTag, "Server is listening on port "+server.Port)

	for {
		// Accept a connection
		conn, err1 := listener.Accept()
		if err1 != nil {
			fmt.Println("Error accepting connection:", err1)
			continue
		}

		// Handle the connection in a new goroutine
		go server.handleConnection(conn)
	}
}
