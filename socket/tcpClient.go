package socket

import (
	"fmt"
	"net"
)

const (
	LoggerTagClient = "TCP Socket Client"
)

type TCPClinet struct {
	Port string
	conn net.Conn
}

func (client *TCPClinet) Connect() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	client.conn = conn
}

func (client *TCPClinet) Close() {
	err := client.conn.Close()
	if err != nil {
		return
	}
}

func (client *TCPClinet) SendData(data []byte) {

	// Send the input to the server
	_, err := client.conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	// Read the response from the server
	response := readData(client.conn)
	if response != nil {
		fmt.Printf("Server response: %s\n", string(response[:]))
	}

}
