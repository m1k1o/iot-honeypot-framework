package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"m1k1o/ioth/internal/config"
	"m1k1o/ioth/internal/utils"
)

type ApiManagerCtx struct {
	logger zerolog.Logger
	conf   *config.API
	cli    *client.Client
}

func New(conf *config.API) *ApiManagerCtx {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	return &ApiManagerCtx{
		logger: log.With().Str("module", "router").Logger(),
		conf:   conf,
		cli:    cli,
	}
}

func (a *ApiManagerCtx) Mount(r chi.Router) {
	// Register
	r.Mount("/nodes", a.nodes())
	r.Mount("/services", a.services())
	r.Mount("/proxies", a.proxies())

	r.Get("/images", func(w http.ResponseWriter, r *http.Request) {
		// List all images available locally
		images, err := a.cli.ImageList(r.Context(), types.ImageListOptions{})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	})

	r.Get("/containers", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve a list of containers
		containers, err := a.cli.ContainerList(r.Context(), types.ContainerListOptions{})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(containers)
	})

	r.Get("/networks", func(w http.ResponseWriter, r *http.Request) {
		// List all networks
		networks, err := a.cli.NetworkList(r.Context(), types.NetworkListOptions{})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(networks)
	})

	// ref.
	// https://godoc.org/github.com/docker/docker/client

	// services
	// [x] func (cli *Client) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error)
	// [ ] func (cli *Client) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error)
	// [x] func (cli *Client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error)
	// [ ] func (cli *Client) ServiceLogs(ctx context.Context, serviceID string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	// [x] func (cli *Client) ServiceRemove(ctx context.Context, serviceID string) error
	// [ ] func (cli *Client) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error)

	// swarm
	// [ ] func (cli *Client) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error)
	// [ ] func (cli *Client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error)
	// [ ] func (cli *Client) SwarmInspect(ctx context.Context) (swarm.Swarm, error)
	// [ ] func (cli *Client) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error
	// [ ] func (cli *Client) SwarmLeave(ctx context.Context, force bool) error
	// [ ] func (cli *Client) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error
	// [ ] func (cli *Client) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, flags swarm.UpdateFlags) error

	// nodes
	// [ ] func (cli *Client) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error)
	// [x] func (cli *Client) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error)
	// [ ] func (cli *Client) NodeRemove(ctx context.Context, nodeID string, options types.NodeRemoveOptions) error
	// [ ] func (cli *Client) NodeUpdate(ctx context.Context, nodeID string, version swarm.Version, node swarm.NodeSpec) error
}

func (a *ApiManagerCtx) imageName(name string) string {
	registry := strings.TrimRight(a.conf.Registry, "/")
	return registry + "/" + name
}
