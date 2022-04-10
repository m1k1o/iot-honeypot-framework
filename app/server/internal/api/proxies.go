package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/go-chi/chi"

	"m1k1o/ioth/internal/utils"
)

type PortConfig struct {
	Protocol swarm.PortConfigProtocol `json:"protocol,omitempty"`
	// TargetPort is the port inside the container
	TargetPort uint32 `json:"target_port,omitempty"`
	// PublishedPort is the port on the swarm hosts
	PublishedPort uint32 `json:"published_port,omitempty"`
}

func PortConfigFromEnv(envs []string) []PortConfig {
	portsEnv := ""
	for _, env := range envs {
		if strings.HasPrefix(env, "PORTS=") {
			portsEnv = strings.TrimPrefix(env, "PORTS=")
			break
		}
	}

	ports := []PortConfig{}
	rows := strings.Split(strings.ReplaceAll(portsEnv, "\r\n", "\n"), "\n")
	for _, row := range rows {
		cols := strings.Fields(row)
		if len(cols) != 3 || strings.HasPrefix(cols[0], "#") {
			continue
		}

		targetPort, _ := strconv.Atoi(cols[1])
		publishedPort, _ := strconv.Atoi(cols[2])
		ports = append(ports, PortConfig{
			Protocol:      swarm.PortConfigProtocol(cols[0]),
			TargetPort:    uint32(targetPort),
			PublishedPort: uint32(publishedPort),
		})
	}

	return ports
}

type ProxySpec struct {
	ID      string       `json:"id,omitempty"`
	Service string       `json:"service,omitempty"`
	Node    string       `json:"node,omitempty"`
	Ports   []PortConfig `json:"ports,omitempty"`
	Running bool         `json:"running"`
}

func (a *ApiManagerCtx) proxies() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		servicesList, err := utils.GetServiceListWithStatus(r.Context(), a.cli)
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		proxies := []ProxySpec{}
		for _, service := range servicesList {
			val, ok := service.Spec.Labels["m1k1o.ioth.type"]
			if !ok || val != "proxy" {
				continue
			}

			proxies = append(proxies, ProxySpec{
				ID:      service.ID,
				Service: service.Spec.Labels["m1k1o.ioth.proxy.service"],
				Node:    service.Spec.Labels["m1k1o.ioth.proxy.node"],
				Ports:   PortConfigFromEnv(service.Spec.TaskTemplate.ContainerSpec.Env),
				Running: service.ServiceStatus.RunningTasks == service.ServiceStatus.DesiredTasks,
			})
		}

		utils.HttpSuccess(w, proxies)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			utils.HttpBadRequest(w, "Please send a request body.")
			return
		}

		var proxy ProxySpec
		if !utils.HttpJsonRequest(w, r, &proxy) {
			return
		}

		// create PortConfig & PortEnv from []PortProxy
		PortConfig := []swarm.PortConfig{}
		buf := bytes.Buffer{}
		obj := make(map[string]bool)
		for _, port := range proxy.Ports {
			PortConfig = append(PortConfig, swarm.PortConfig{
				//Name:          "foo",
				Protocol:      port.Protocol,
				TargetPort:    port.PublishedPort,
				PublishedPort: port.PublishedPort,
				PublishMode:   swarm.PortConfigPublishModeHost,
			})

			entry := fmt.Sprintf("%s %d %d\n", port.Protocol, port.TargetPort, port.PublishedPort)
			if _, ok := obj[entry]; !ok {
				obj[entry] = true
				buf.WriteString(entry)
			}
		}
		PortEnv := buf.String()

		a.logger.Info().
			Str("Service", proxy.Service).
			Str("Node", proxy.Node).
			Str("PortEnv", PortEnv).
			Msgf("create proxy")

		// labels
		Labels := make(map[string]string)
		Labels["m1k1o.ioth.type"] = "proxy"
		Labels["m1k1o.ioth.proxy.service"] = proxy.Service
		Labels["m1k1o.ioth.proxy.node"] = proxy.Node

		Annotations := swarm.Annotations{
			//Name:  "foo-bar",
			Labels: Labels,
		}

		// https://godoc.org/github.com/docker/docker/api/types/swarm
		TaskTemplate := swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: a.imageName("ioth-docker-proxy"),
				//Labels          map[string]string,
				//Command         []string,
				//Args            []string,
				//Hostname        string,
				Env: []string{
					"HOST=" + proxy.Service,
					"PORTS=" + PortEnv,
				},
				//Dir             string,
				//User            string,
				//Groups          []string,
				//Privileges      *Privileges,
				//Init            *bool,
				//StopSignal      string,
				//TTY             bool,
				//OpenStdin       bool,
				//ReadOnly        bool,
				//Mounts          []mount.Mount,
				//StopGracePeriod *time.Duration,
				//Healthcheck     *container.HealthConfig,
				//Hosts           []string,
				//DNSConfig       *DNSConfig,
				//Secrets         []*SecretReference,
				//Configs         []*ConfigReference,
				//Isolation       container.Isolation,
				//Sysctls         map[string]string,
				//CapabilityAdd   []string,
				//CapabilityDrop  []string,
				//Ulimits         []*units.Ulimit,
			},
			Resources: &swarm.ResourceRequirements{},
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyConditionAny,
			},
			Placement: &swarm.Placement{
				Constraints: []string{
					"node.hostname==" + proxy.Node,
				},
			},
			LogDriver: &swarm.Driver{
				Name: "fluentd",
				Options: map[string]string{
					"tag":                   "proxy." + proxy.Service + "." + proxy.Node,
					"fluentd-address":       a.conf.Fluentd,
					"fluentd-async-connect": "true",
				},
			},
		}

		Mode := swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: func(i uint64) *uint64 { return &i }(1),
			},
			//Global: &swarm.GlobalService{},
		}

		Networks := []swarm.NetworkAttachmentConfig{
			{
				Target: a.conf.OverlayNetwork,
			},
			{
				Target: "bridge",
			},
		}

		EndpointSpec := &swarm.EndpointSpec{
			Ports: PortConfig,
		}

		// ---

		service := swarm.ServiceSpec{
			Annotations:  Annotations,
			TaskTemplate: TaskTemplate,
			Mode:         Mode,
			UpdateConfig: &swarm.UpdateConfig{},
			Networks:     Networks,
			EndpointSpec: EndpointSpec,
		}

		ServiceCreateResponse, err := a.cli.ServiceCreate(r.Context(), service, types.ServiceCreateOptions{})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		//ServiceCreateResponse.ID
		utils.HttpSuccess(w, ServiceCreateResponse)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		service, _, err := a.cli.ServiceInspectWithRaw(r.Context(), id, types.ServiceInspectOptions{})
		if err != nil {
			utils.HttpBadRequest(w, err)
			return
		}

		val, ok := service.Spec.Labels["m1k1o.ioth.type"]
		if !ok || val != "proxy" {
			utils.HttpBadRequest(w, "Swarm service is not ioth proxy.")
			return
		}

		if err := a.cli.ServiceRemove(r.Context(), id); err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		utils.HttpSuccess(w)
	})

	return r
}
