package socket

import (
	"fmt"
	"net"
)

func readData(conn net.Conn) []byte {

	output := make([]byte, 0)
	for {
		// Read data from the connection
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return nil
		}

		// Print the received data
		fmt.Printf("Received: %s", buffer[:n])

		output = append(output, buffer[:n]...)

		if n < 1024 {
			break
		}

	}

	return output
}
