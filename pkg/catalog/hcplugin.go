package catalog

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/zeebo/errs/v2"
	"google.golang.org/grpc"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/openkcm/plugin-sdk/internal/consts"
)

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type HCRPCPlugin struct {
	// HCRPCPlugin must implement the Plugin interface
	goplugin.NetRPCUnsupportedPlugin

	config PluginConfig
}

var _ goplugin.GRPCPlugin = (*HCRPCPlugin)(nil)

func (p *HCRPCPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	return errors.New("not implemented host side")
}

func (p *HCRPCPlugin) GRPCClient(ctx context.Context, b *goplugin.GRPCBroker, c *grpc.ClientConn) (any, error) {
	// Manually start up the server via b.Accept since b.AcceptAndServe does
	// some logging we don't care for. Although b.AcceptAndServe is currently
	// the only way to feed the TLS config to the brokered connection, AutoMTLS
	// does not work yet anyway, so it is a moot point.
	listener, err := b.Accept(consts.HostServiceProviderID)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	server := newHostServer(p.config.Logger, p.config.Name)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			p.config.Logger.ErrorContext(ctx, "Host services server failed", "error", err)
			if err := c.Close(); err != nil {
				p.config.Logger.ErrorContext(ctx, "Failed to close client connection", "error", err)
			}
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if !gracefulStopWithTimeout(server, time.Minute) {
			p.config.Logger.WarnContext(ctx, "Forced timed-out host service server to stop")
		}
	}()

	return &HCPlugin{
		conn:    c,
		closers: groupCloserFuncs(cancel, wg.Wait),
	}, nil
}

type HCPlugin struct {
	conn grpc.ClientConnInterface

	closers closerGroup
}
