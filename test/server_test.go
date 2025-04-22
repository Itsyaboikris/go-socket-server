// test/server_test.go
package test

import (
	"testing"

	"github.com/itsyaboikris/go_socket_server/client"
)

func TestServerResponse(t *testing.T) {
    // Setup test server
    ts := NewTestServer(":9001")
    err := ts.Start()
    if err != nil {
        t.Fatalf("Failed to start server: %v", err)
    }
    defer ts.Stop()

    // Test cases
    tests := []struct {
        name    string
        message string
        want    string
    }{
        {
            name:    "Simple message",
            message: "Hello Server!",
            want:    "Hello Server!\n",
        },
        {
            name:    "Empty message",
            message: "",
            want:    "\n",
        },
    }

    // Run test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            response, err := client.ConnectAndSend(ts.address, tt.message)
            if err != nil {
                t.Fatalf("ConnectAndSend failed: %v", err)
            }

            if response != tt.want {
                t.Errorf("got %q, want %q", response, tt.want)
            }
        })
    }
}
