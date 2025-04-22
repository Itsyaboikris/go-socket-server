// test/load_test.go
package test

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	batchSize         = 100   // Number of concurrent connections per batch
	totalConnections  = 10000 // Total number of connections to test
	connectionTimeout = 10 * time.Second
)

func BenchmarkConcurrentConnections(b *testing.B) {
	ts := NewTestServer(":9002")
	err := ts.Start()
	if err != nil {
		b.Fatalf("Failed to start server: %v", err)
	}
	defer ts.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Process connections in batches
		for start := 0; start < totalConnections; start += batchSize {
			end := start + batchSize
			if end > totalConnections {
				end = totalConnections
			}

			var wg sync.WaitGroup
			errCh := make(chan error, batchSize)

			// Launch connections for this batch
			for j := start; j < end; j++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()

					// Set timeout for the connection
					dialer := net.Dialer{Timeout: connectionTimeout}
					conn, err := dialer.Dial("tcp", ts.address)
					if err != nil {
						errCh <- fmt.Errorf("connection %d failed: %v", id, err)
						return
					}
					defer conn.Close()

					// Set deadlines for read/write operations
					conn.SetDeadline(time.Now().Add(connectionTimeout))

					// Send message
					message := "Hello Server!"
					if _, err := conn.Write([]byte(message)); err != nil {
						errCh <- fmt.Errorf("write failed on connection %d: %v", id, err)
						return
					}

					// Read response
					buffer := make([]byte, 1024)
					expected := len(message)
					totalRead := 0
					for totalRead < expected {
						n, err := conn.Read(buffer[totalRead:])
						if err != nil {
							errCh <- fmt.Errorf("read failed on connection %d: %v", id, err)
							return
						}
						totalRead += n
					}

					response := string(buffer[:totalRead])
					if response != message {
						errCh <- fmt.Errorf("connection %d: expected %q, got %q", id, message, response)
					}
				}(j)
			}

			// Wait for current batch to complete
			wg.Wait()
			close(errCh)

			// Check for any errors in this batch
			for err := range errCh {
				b.Fatal(err)
			}

			// Add small delay between batches to prevent overwhelming the server
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func BenchmarkHighThroughputConnections(b *testing.B) {
	ts := NewTestServer(":9003")
	err := ts.Start()
	if err != nil {
		b.Fatalf("Failed to start server: %v", err)
	}
	defer ts.Stop()

	b.ResetTimer()

	var failed atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, err := net.DialTimeout("tcp", ts.address, connectionTimeout)
			if err != nil {
				b.Errorf("dial error: %v", err)
				return
			}

			conn.SetDeadline(time.Now().Add(connectionTimeout))
			message := "Hello Server!"
			if _, err := conn.Write([]byte(message)); err != nil {
				b.Errorf("write error: %v", err)
				conn.Close()
				return
			}

			buffer := make([]byte, 1024)
			expected := len(message)
			totalRead := 0
			for totalRead < expected {
				n, err := conn.Read(buffer[totalRead:])
				if err != nil {
					b.Errorf("read error: %v", err)
					conn.Close()
					return
				}
				totalRead += n
			}

			response := string(buffer[:totalRead])
			if response != message {
				b.Errorf("expected %q, got %q", message, response)
				failed.Add(1)
			}

			// Optional: Sleep before closing the connection
			time.Sleep(10 * time.Millisecond)

			// Log any failures
			if failed.Load() > 0 {
				b.Logf("Failed connections: %d", failed.Load())
			}
			conn.Close()
		}
	})

	// Optionally log total failures at the end
	b.Logf("Total failed connections: %d", failed.Load())
}
