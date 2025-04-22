package server

import (
	"context"
	"io"
	"log"
	"net"
	"time"
)

const (
	maxConcurrentConnections = 10000
	readTimeout              = 10 * time.Second
	writeTimeout             = 10 * time.Second
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		conn.SetReadDeadline(time.Now().Add(readTimeout))
		message, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error reading from connection: %v", err)
			continue
		}

		conn.SetWriteDeadline(time.Now().Add(writeTimeout))
		_, err = conn.Write(buffer[:message])
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}

func Start(ctx context.Context, address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", address, err)
	}
	defer listener.Close()
	log.Printf("Server listening on %s", address)

	// Semaphore: limit to maxConcurrentConnections
	sem := make(chan struct{}, maxConcurrentConnections)

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		listener.Close()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			listener.(*net.TCPListener).SetDeadline(time.Now().Add(10 * time.Second))
			conn, err := listener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if ctx.Err() != nil {
					return
				}
				log.Printf("Failed to accept connection: %v", err)
				continue
			}

			select {
			case sem <- struct{}{}:
				go func() {
					defer func() { <-sem }() // Release slot
					handleConnection(conn)
				}()
			default:
				log.Println("Too many connections – rejecting new one")
				time.Sleep(100 * time.Millisecond)
				conn.Close()
			}
		}
	}
}
