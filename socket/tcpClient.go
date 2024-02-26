package socket

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"net"
)

type TCPClient struct {
	Port       string
	Host       string
	conn       net.Conn
	WaitForAck bool
}

func CreateTCPClient(clientConfig TCPClientConfig) *TCPClient {
	return &TCPClient{Host: clientConfig.Host, Port: clientConfig.Port, WaitForAck: clientConfig.WaitForAck}
}

func (client *TCPClient) Connect() bool {
	// Connect to the server
	conn, err := net.Dial("tcp", client.Host+":"+client.Port)
	if err != nil {
		zkLogger.Error(LoggerTagSocket, "Error connecting:", err)
		return false
	}
	client.conn = conn
	return true
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
		zkLogger.Error(LoggerTagSocket, "Error sending data:", err)
		return
	}

	// Read the response from the server
	if client.WaitForAck {
		response := readData(client.conn)
		if response != nil {
			zkLogger.DebugF(LoggerTagSocket, "Received on client: %s\n", string(response[:]))
		}
	}

}
