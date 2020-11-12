package mortar

import (
	"github.com/go-masonry/bprometheus"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/monitor"
	mortarProject "github.com/go-masonry/mortar/mortar"
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

// PrometheusFxOption registers prometheus
func PrometheusFxOption() fx.Option {
	return fx.Options(
		providers.MonitorFxOption(),
		providers.MonitorGRPCInterceptorFxOption(),
		bprometheus.PrometheusInternalHandlerFxOption(),
		fx.Provide(PrometheusBuilder),
	)
}

// PrometheusBuilder returns a monitor.Builder that is implemented by Prometheus
func PrometheusBuilder(cfg cfg.Config) monitor.Builder {
	name := cfg.Get(mortarProject.Name).String()
	return bprometheus.Builder().SetNamespace(name)
}
