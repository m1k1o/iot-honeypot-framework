package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"m1k1o/ioth"
	"m1k1o/ioth/cmd"
	"m1k1o/ioth/internal/utils"
)

func main() {
	fmt.Print(utils.Colorf(ioth.Header, "server", ioth.Service.Version))
	if err := cmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("failed to execute command")
	}
}
