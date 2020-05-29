package extension

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Bar interface {
	Bars() []string
}

type BarRPC struct{ client *rpc.Client }

func (g *BarRPC) Bars() (resp []string) {
	err := g.client.Call("Plugin.Bars", map[string]interface{}{}, &resp)
	if err != nil {
		panic(err)
	}
	return
}

type BarRPCServer struct {
	Impl Bar
}

func (s *BarRPCServer) Bars(args map[string]interface{}, resp *[]string) (err error) {
	*resp = s.Impl.Bars()
	return
}

type BarPlugin struct {
	// Impl Injection
	Impl Bar
}

func (p *BarPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &BarRPCServer{Impl: p.Impl}, nil
}

func (BarPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &BarRPC{client: c}, nil
}
