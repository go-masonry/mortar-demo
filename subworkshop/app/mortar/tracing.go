package mortar

import (
	"context"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"

	"github.com/go-masonry/bjaeger"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
)

func TracerFxOption() fx.Option {
	return fx.Provide(JaegerBuilder)
}

// This constructor assumes you have JAEGER environment variables set
//
// https://github.com/jaegertracing/jaeger-client-go#environment-variables
//
// Once built it will register Lifecycle hooks (connect on start, close on stop)
func JaegerBuilder(lc fx.Lifecycle, config cfg.Config, logger log.Logger) (opentracing.Tracer, error) {
	openTracer, err := bjaeger.Builder().
		SetServiceName(config.Get(confkeys.ApplicationName).String()).
		AddOptions(bjaeger.BricksLoggerOption(logger)). // verbose logging,
		Build()
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return openTracer.Connect(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return openTracer.Close(ctx)
		},
	})
	return openTracer.Tracer(), nil
}
