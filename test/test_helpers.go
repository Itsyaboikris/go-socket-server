// test/test_helpers.go
package test

import (
	"context"
	"sync"
	"time"

	"github.com/itsyaboikris/go_socket_server/server"
)

type TestServer struct {
	address string
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewTestServer(address string) *TestServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TestServer{
		address: address,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (ts *TestServer) Start() error {
	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		server.Start(ts.ctx, ts.address)
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (ts *TestServer) Stop() {
	ts.cancel()
	ts.wg.Wait()
}
