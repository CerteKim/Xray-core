package web

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/xtls/xray-core/app/web/config"
	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/signal/done"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/outbound"
)

type WebHandler struct {
	ohm    outbound.Manager
	tag    string
	api    Api
	pprof  bool
	static []Static
}

type Api struct {
	address string
	port    uint32
}

type Static struct {
	filePath string
	uri      string
}

// New
func NewWebHandler(ctx context.Context, config *config.Config) (*WebHandler, error) {
	c := &WebHandler{
		tag: config.Tag,
		api: Api{
			address: config.Api.Address,
			port:    config.Api.Port,
		},
		pprof: config.Pprof,
	}

	for _, s := range config.Static {
		c.static = append(c.static, Static{
			filePath: s.FilePath,
			uri:      s.Uri,
		})
	}

	common.Must(core.RequireFeatures(ctx, func(om outbound.Manager) {
		c.ohm = om
	}))
	return c, nil
}

func (r *WebHandler) Type() interface{} {
	return (*WebHandler)(nil)
}

func (r *WebHandler) Start() error {
	listener := &OutboundListener{
		buffer: make(chan net.Conn, 4),
		done:   done.New(),
	}

	h2s := &http2.Server{}

	go func() {
		if err := http.Serve(listener, h2c.NewHandler(Default(r), h2s)); err != nil {
			newError("failed to start Web server").Base(err).AtError().WriteToLog()
		}
	}()

	if err := r.ohm.RemoveHandler(context.Background(), r.tag); err != nil {
		newError("failed to remove existing handler").WriteToLog()
	}

	return r.ohm.AddHandler(context.Background(), &Outbound{
		tag:      r.tag,
		listener: listener,
	})
}

func (r *WebHandler) Close() error {
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*config.Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		return NewWebHandler(ctx, cfg.(*config.Config))
	}))
}
