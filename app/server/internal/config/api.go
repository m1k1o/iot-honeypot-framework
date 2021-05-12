package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type API struct {
	Registry       string
	Fluentd        string
	OverlayNetwork string
}

func (API) Init(cmd *cobra.Command) error {
	cmd.PersistentFlags().String("registry", "ioth-registry:5000", "address of the internal registry")
	if err := viper.BindPFlag("registry", cmd.PersistentFlags().Lookup("registry")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("fluentd", "10.8.0.1:24224", "address of the fluentd service")
	if err := viper.BindPFlag("fluentd", cmd.PersistentFlags().Lookup("fluentd")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("overlay_network", "ioth_overlay_app", "name of the attachable overlay network")
	if err := viper.BindPFlag("overlay_network", cmd.PersistentFlags().Lookup("overlay_network")); err != nil {
		return err
	}

	return nil
}

func (s *API) Set() {
	s.Registry = viper.GetString("registry")
	s.Fluentd = viper.GetString("fluentd")
	s.OverlayNetwork = viper.GetString("overlay_network")
}
