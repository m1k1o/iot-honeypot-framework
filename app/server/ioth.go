package ioth

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"m1k1o/ioth/internal/api"
	"m1k1o/ioth/internal/config"
	"m1k1o/ioth/internal/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const Header = `&34
    ____    ______   __                                       __ 
   /  _/___/_  __/  / /_  ____  ____  ___  __  ______  ____  / /_
   / // __ \/ /    / __ \/ __ \/ __ \/ _ \/ / / / __ \/ __ \/ __/
 _/ // /_/ / /    / / / / /_/ / / / /  __/ /_/ / /_/ / /_/ / /_  
/___/\____/_/    /_/ /_/\____/_/ /_/\___/\__, / .___/\____/\__/  
                                        /____/_/                 
&1&37                     by m1k1o                    &33%s v%s&0
`

var (
	//
	buildDate = "dev"
	//
	gitCommit = "dev"
	//
	gitBranch = "dev"

	// Major version when you make incompatible API changes,
	major = "2"
	// Minor version when you add functionality in a backwards-compatible manner, and
	minor = "0"
	// Patch version when you make backwards-compatible bug fixeioth.
	patch = "0"
)

var Service *Ioth

func init() {
	Service = &Ioth{
		Version: &Version{
			Major:     major,
			Minor:     minor,
			Patch:     patch,
			GitCommit: gitCommit,
			GitBranch: gitBranch,
			BuildDate: buildDate,
			GoVersion: runtime.Version(),
			Compiler:  runtime.Compiler,
			Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		},
		Configs: &Configs{
			Root:   &config.Root{},
			Server: &config.Server{},
			API:    &config.API{},
		},
	}
}

type Version struct {
	Major     string
	Minor     string
	Patch     string
	GitCommit string
	GitBranch string
	BuildDate string
	GoVersion string
	Compiler  string
	Platform  string
}

func (i *Version) String() string {
	return fmt.Sprintf("%s.%s.%s %s", i.Major, i.Minor, i.Patch, i.GitCommit)
}

func (i *Version) Details() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		fmt.Sprintf("Version %s.%s.%s", i.Major, i.Minor, i.Patch),
		fmt.Sprintf("GitCommit %s", i.GitCommit),
		fmt.Sprintf("GitBranch %s", i.GitBranch),
		fmt.Sprintf("BuildDate %s", i.BuildDate),
		fmt.Sprintf("GoVersion %s", i.GoVersion),
		fmt.Sprintf("Compiler %s", i.Compiler),
		fmt.Sprintf("Platform %s", i.Platform),
	)
}

type Configs struct {
	Root   *config.Root
	Server *config.Server
	API    *config.API
}

type Ioth struct {
	Version *Version
	Configs *Configs

	logger      zerolog.Logger
	apiManager  *api.ApiManagerCtx
	httpManager *http.HttpManagerCtx
}

func (ioth *Ioth) Preflight() {
	ioth.logger = log.With().Str("service", "ioth").Logger()
}

func (ioth *Ioth) Start() {
	ioth.apiManager = api.New(
		ioth.Configs.API,
	)

	ioth.httpManager = http.New(
		ioth.apiManager,
		ioth.Configs.Server,
	)
	ioth.httpManager.Start()
}

func (ioth *Ioth) Shutdown() {
	if err := ioth.httpManager.Shutdown(); err != nil {
		ioth.logger.Err(err).Msg("http manager shutdown with an error")
	} else {
		ioth.logger.Debug().Msg("http manager shutdown")
	}
}

func (ioth *Ioth) ServeCommand(cmd *cobra.Command, args []string) {
	ioth.logger.Info().Msg("starting ioth server")
	ioth.Start()
	ioth.logger.Info().Msg("ioth ready")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	ioth.logger.Warn().Msgf("received %s, attempting graceful shutdown: \n", sig)
	ioth.Shutdown()
	ioth.logger.Info().Msg("shutdown complete")
}
