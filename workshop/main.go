package main

import (
	"github.com/alecthomas/kong"
	"github.com/go-masonry/mortar-demo/workshop/app/mortar"
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

var CLI struct {
	Config struct {
		Path            string   `arg:"" required:"" help:"Path to config file." type:"existingfile"`
		AdditionalFiles []string `optional:"" help:"Additional configuration files to merge, comma separated" type:"existingfile"`
	} `cmd:"" help:"Path to config file."`
}

func main() {
	ctx := kong.Parse(&CLI, kong.UsageOnError())
	switch cmd := ctx.Command(); cmd {
	case "config <path>":
		app := createApplication(CLI.Config.Path, CLI.Config.AdditionalFiles)
		app.Run()
	default:
		ctx.Fatalf("unknown option %s", cmd)
	}
}

func createApplication(configFilePath string, additionalFiles []string) *fx.App {
	return fx.New(
		fx.NopLogger,                                             // remove fx debug
		mortar.ViperFxOption(configFilePath, additionalFiles...), // Configuration map
		mortar.LoggerFxOption(),                                  // Logger
		mortar.TracerFxOption(),                                  // Jaeger tracing
		mortar.PrometheusFxOption(),                              // Prometheus
		mortar.HttpClientFxOptions(),
		mortar.HttpServerFxOptions(),
		mortar.InternalHttpHandlersFxOptions(),
		// Tutorial service dependencies
		mortar.WorkshopAPIsAndOtherDependenciesFxOption(), // register tutorial APIs
		// This one invokes all the above
		providers.BuildMortarWebServiceFxOption(), // http server invoker
	)
}
