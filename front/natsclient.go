package front

import (
	"context"
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/nats-io/nats.go"
	"net"
	"nhooyr.io/websocket"
	"time"
)

var _ nats.CustomDialer = (*WasmNatsConnectionWrapper)(nil) // Verify the implementation

func (cw WasmNatsConnectionWrapper) Dial(network, address string) (net.Conn, error) {
	// we actually do not care about the adress given here
	app.Logf("Got Request for Network: %q / Address: %q", network, address)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	app.Logf("Dialing Address: ws://%s", address)
	c, _, err := websocket.Dial(ctx, "ws://"+address, nil)
	if err != nil {
		app.Logf("websocket.Dial failed %#v", err)
		return nil, fmt.Errorf("websocket.Dial failed %w", err)
	}
	cw.ws = c
	nconn := websocket.NetConn(context.Background(), c, websocket.MessageBinary)
	return nconn, nil
}

func (cw WasmNatsConnectionWrapper) SkipTLSHandshake() bool {
	return true
}
