package example

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// PingPonger is the interface that we're exposing as a plugin.
type PingPonger interface {
	Ping() (string, error)
}

// Here is an implementation that talks over RPC
type PingPongerRPC struct{ client *rpc.Client }

func (g *PingPongerRPC) Ping() (resp string, err error) {
	err = g.client.Call("Plugin.Ping", map[string]interface{}{}, &resp)
	return
}

// PingPongerRPCServer ping pong server
type PingPongerRPCServer struct {
	// This is the real implementation
	Impl PingPonger
}

func (s *PingPongerRPCServer) Ping(args map[string]interface{}, resp *string) (err error) {
	*resp, err = s.Impl.Ping()
	return
}

type PingPongerPlugin struct {
	// Impl Injection
	Impl PingPonger
}

func (p *PingPongerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PingPongerRPCServer{Impl: p.Impl}, nil
}

func (PingPongerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PingPongerRPC{client: c}, nil
}
