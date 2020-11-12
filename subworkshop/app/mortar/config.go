package mortar

import (
	"github.com/go-masonry/bviper"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"go.uber.org/fx"
)

func ViperFxOption(configFilePath string, additionalFilePaths ...string) fx.Option {
	return fx.Provide(func() (cfg.Config, error) {
		builder := bviper.Builder().SetConfigFile(configFilePath)
		for _, extraFile := range additionalFilePaths {
			builder = builder.AddExtraConfigFile(extraFile)
		}
		return builder.Build()
	})
}
