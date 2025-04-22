package client

import (
	"bufio"
	"fmt"
	"net"
)

func ConnectAndSend(address string, message string) (string, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	fmt.Fprintf(conn, message+"\n")
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return response, nil
}
