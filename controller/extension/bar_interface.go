package extension

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"

	"github.com/danielfbm/plugin-example/controller/extension/bar"
	"google.golang.org/grpc"
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

type BarGRPC struct {
	plugin.Plugin
	Impl Bar
}

func (p *BarGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bar.RegisterBarServer(s, &BarGRPCServer{Impl: p.Impl})
	return nil
}

func (p *BarGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &BarGRPCClient{client: bar.NewBarClient(c)}, nil
}

type BarGRPCClient struct {
	client bar.BarClient
}

func (c *BarGRPCClient) Bars() (resp []string) {
	response, _ := c.client.Bars(context.Background(), &bar.Empty{})
	if response != nil {
		resp = response.Value
	}
	return
}

type BarGRPCServer struct {
	Impl Bar
}

func (s *BarGRPCServer) Bars(ctx context.Context, _ *bar.Empty) (*bar.BarsResponse, error) {
	resp := &bar.BarsResponse{}
	resp.Value = s.Impl.Bars()
	return resp, nil
}
