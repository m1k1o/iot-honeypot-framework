package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"m1k1o/ioth"
	"m1k1o/ioth/internal/config"
)

func init() {
	command := &cobra.Command{
		Use:   "serve",
		Short: "serve ioth manager server",
		Long:  `serve ioth manager server`,
		Run:   ioth.Service.ServeCommand,
	}

	configs := []config.Config{
		ioth.Service.Configs.Server,
		ioth.Service.Configs.API,
	}

	cobra.OnInitialize(func() {
		for _, cfg := range configs {
			cfg.Set()
		}
		ioth.Service.Preflight()
	})

	for _, cfg := range configs {
		if err := cfg.Init(command); err != nil {
			log.Panic().Err(err).Msg("unable to run serve command")
		}
	}

	root.AddCommand(command)
}
