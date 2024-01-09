package socket

import (
	zklogger "github.com/zerok-ai/zk-utils-go/logs"
	"net"
)

const (
	LoggerTagSocket = "TCP Socket"
)

func readData(conn net.Conn) []byte {

	output := make([]byte, 0)
	for {
		// Read data from the connection
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			zklogger.Error(LoggerTagSocket, "Error reading:", err)
			return nil
		}

		output = append(output, buffer[:n]...)

		if n < 1024 {
			break
		}

	}

	return output
}
