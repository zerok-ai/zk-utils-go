package socket

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"net"
)

const (
	LoggerTagSocket = "TCP Socket"
	bufferSize      = 1024
)

type TCPServerConfig struct {
	Port    string `yaml:"port" env:"TCP_SERVER_PORT" env-description:"Server port" env-default:"6473"`
	SendAck bool   `yaml:"sendAck" env:"TCP_SERVER_SEND_ACK" env-description:"Server to acknowledge the message to clinet" env-default:"false"`
}

type TCPClientConfig struct {
	Host       string `yaml:"host" env:"TCP_CLIENT_HOST" env-description:"Client host" env-default:"127.0.0.1"`
	Port       string `yaml:"port" env:"TCP_CLIENT_PORT" env-description:"Client port" env-default:"6473"`
	WaitForAck bool   `yaml:"waitForAck" env:"TCP_SERVER_WAIT_FOR_ACK" env-description:"Server to acknowledge the message to client" env-default:"false"`
}

func readData(conn net.Conn) []byte {

	output := make([]byte, 0)
	for {
		// Read data from the connection
		buffer := make([]byte, bufferSize)
		n, err := conn.Read(buffer)
		if err != nil {
			zkLogger.Error(LoggerTagSocket, "Error reading:", err)
			return nil
		}

		output = append(output, buffer[:n]...)

		if n < bufferSize {
			break
		}

	}

	return output
}

func writeData(conn net.Conn, data []byte) bool {

	_, err := conn.Write(data)
	if err != nil {
		zkLogger.Error(LoggerTagSocket, "Error writing:", err)
		return false
	}

	return true
}
