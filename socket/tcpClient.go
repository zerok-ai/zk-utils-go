package socket

import (
	"fmt"
	"net"
)

const (
	LoggerTagClient = "TCP Socket Client"
)

type TCPClient struct {
	Port string
	Host string
	conn net.Conn
}

func (client *TCPClient) Connect() {
	// Connect to the server
	conn, err := net.Dial("tcp", client.Host+":"+client.Port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	client.conn = conn
}

func (client *TCPClient) Close() {
	if client.conn == nil {
		return
	}
	err := client.conn.Close()
	if err != nil {
		return
	}
}

func (client *TCPClient) SendData(data []byte) {

	if client.conn == nil {
		return
	}

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
