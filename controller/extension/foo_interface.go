package extension

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Foo interface {
	Foos() (string, error)
}

type FooRPC struct{ client *rpc.Client }

func (g *FooRPC) Foos() (resp string, err error) {
	err = g.client.Call("Plugin.Foos", map[string]interface{}{}, &resp)
	return
}

type FooRPCServer struct {
	// This is the real implementation
	Impl Foo
}

func (s *FooRPCServer) Foos(args map[string]interface{}, resp *string) (err error) {
	*resp, err = s.Impl.Foos()
	return
}

type FooPlugin struct {
	// Impl Injection
	Impl Foo
}

func (p *FooPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FooRPCServer{Impl: p.Impl}, nil
}

func (FooPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FooRPC{client: c}, nil
}
