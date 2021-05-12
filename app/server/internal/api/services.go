package api

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi"

	"m1k1o/ioth/internal/utils"
)

type ServiceStatus struct {
	Running uint64 `json:"running"`
	Desired uint64 `json:"desired"`
}

type ContainerSpec struct {
	Image    string   `json:"image,omitempty"`
	Cmd      []string `json:"cmd,omitempty"`
	Args     []string `json:"args,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	Env      []string `json:"env,omitempty"`
	Dir      string   `json:"dir,omitempty"`
}

type ServiceSpec struct {
	ID            string        `json:"id,omitempty"`
	Name          string        `json:"name,omitempty"`
	ContainerSpec ContainerSpec `json:"container_spec,omitempty"`
	Replicas      uint64        `json:"replicas,omitempty"`
	Status        ServiceStatus `json:"status,omitempty"`
}

func (a *ApiManagerCtx) services() *chi.Mux {
	r := chi.NewRouter()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		servicesList, err := utils.GetServiceListWithStatus(context.Background(), cli)
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		services := []ServiceSpec{}
		for _, service := range servicesList {
			val, ok := service.Spec.Labels["m1k1o.ioth.type"]
			if !ok || val != "service" {
				continue
			}

			services = append(services, ServiceSpec{
				ID:   service.ID,
				Name: service.Spec.Name,
				ContainerSpec: ContainerSpec{
					Image: service.Spec.TaskTemplate.ContainerSpec.Image,
					Env:   service.Spec.TaskTemplate.ContainerSpec.Env,
				},
				Replicas: *service.Spec.Mode.Replicated.Replicas,
				Status: ServiceStatus{
					Running: service.ServiceStatus.RunningTasks,
					Desired: service.ServiceStatus.DesiredTasks,
				},
			})
		}

		utils.HttpSuccess(w, services)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			utils.HttpBadRequest(w, "Please send a request body.")
			return
		}

		var serviceSpec ServiceSpec
		if !utils.HttpJsonRequest(w, r, &serviceSpec) {
			return
		}

		Labels := make(map[string]string)
		Labels["m1k1o.ioth.type"] = "service"

		Annotations := swarm.Annotations{
			Name:   serviceSpec.Name,
			Labels: Labels,
		}

		TaskTemplate := swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: serviceSpec.ContainerSpec.Image,
				//Labels          map[string]string,
				Command:  serviceSpec.ContainerSpec.Cmd,
				Args:     serviceSpec.ContainerSpec.Args,
				Hostname: serviceSpec.ContainerSpec.Hostname,
				Env:      serviceSpec.ContainerSpec.Env,
				Dir:      serviceSpec.ContainerSpec.Dir,
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
			Resources: &swarm.ResourceRequirements{
				//Limits:       &swarm.Limit{
				//	NanoCPUs:    int64(2 * 1e6),
				//	MemoryBytes: int64(4 * 4e6), // at least 4MB
				//	//Pids         int64,
				//},
				//Reservations: &swarm.Resources{
				//	NanoCPUs:    int64(4 * 1e6),
				//	MemoryBytes: int64(8 * 1e6), // at least 4MB
				//	//GenericResources []swarm.GenericResource{
				//	//	{
				//	//		NamedResourceSpec: swarm.NamedGenericResource{
				//	//			Kind  string,
				//	//			Value string,
				//	//		},
				//	//	},
				//	//	{
				//	//		DiscreteResourceSpec: swarm.DiscreteGenericResource{
				//	//			Kind  string,
				//	//			Value int64,
				//	//		},
				//	//	},
				//	//},
				//},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyConditionAny, //or swarm.RestartPolicyConditionNone | swarm.RestartPolicyConditionOnFailure
				//Delay:       *time.Duration,
				//MaxAttempts: *uint64,
				//Window:      *time.Duration,
			},
			Placement: &swarm.Placement{
				//Constraints: []string{"foo", "bar"}, // constraint expected one operator from ==, !=
				//Preferences: []swarm.PlacementPreference{
				//	{
				//		Spread: swarm.SpreadOver{
				//			SpreadDescriptor string,
				//		},
				//	},
				//},
				//MaxReplicas: uint64,
				//Platforms: []swarm.Platform{
				//	{
				//		Architecture: string,
				//		OS:           string,
				//	},
				//}
			},
			LogDriver: &swarm.Driver{
				Name: "fluentd",
				Options: map[string]string{
					"tag":                   "service." + serviceSpec.Name,
					"fluentd-address":       a.conf.Fluentd,
					"fluentd-async-connect": "true",
				},
			},
		}

		Mode := swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &serviceSpec.Replicas,
			},
			//Global:     &swarm.GlobalService{},
			//ReplicatedJob: *ReplicatedJob,
			//GlobalJob:     *GlobalJob,
		}

		Networks := []swarm.NetworkAttachmentConfig{
			{
				Target: a.conf.OverlayNetwork,
				//Aliases: []string{"alias1", "alias2"},
				//DriverOpts: map[string]string,
			},
			//{
			//	Target:  "target_network2",
			//	Aliases: []string{"alias1", "alias2"},
			//	DriverOpts: map[string]string,
			//},
		}

		EndpointSpec := &swarm.EndpointSpec{
			//Mode: swarm.ResolutionModeVIP, //or swarm.ResolutionModeDNSRR
			//Ports: []swarm.PortConfig{
			//	{
			//		Name:          "foo",
			//		Protocol:      swarm.PortConfigProtocolTCP, //or swarm.PortConfigProtocolUDP
			//		TargetPort:    80,
			//		PublishedPort: 8081,
			//		PublishMode:   swarm.PortConfigPublishModeIngress, // swarm.PortConfigPublishModeHost
			//	},
			//},
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

		ServiceCreateResponse, err := cli.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		//ServiceCreateResponse.ID
		utils.HttpSuccess(w, ServiceCreateResponse)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		service, _, err := cli.ServiceInspectWithRaw(context.Background(), id, types.ServiceInspectOptions{})
		if err != nil {
			utils.HttpBadRequest(w, err)
			return
		}

		val, ok := service.Spec.Labels["m1k1o.ioth.type"]
		if !ok || val != "service" {
			utils.HttpBadRequest(w, "Swarm service is not ioth service.")
			return
		}

		if err := cli.ServiceRemove(context.Background(), id); err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		utils.HttpSuccess(w)
	})

	return r
}

/*
[
    {
        "ID": "532bsiap7z3c8l1n4ffz08jcy",
        "Version": {
            "Index": 73
        },
        "CreatedAt": "2020-10-28T15:44:41.876320055Z",
        "UpdatedAt": "2020-10-28T15:44:41.87972199Z",
        "Spec": {
            "Name": "registry",
            "Labels": {},
            "TaskTemplate": {
                "ContainerSpec": {
                    "Image": "registry:2@sha256:8be26f81ffea54106bae012c6f349df70f4d5e7e2ec01b143c46e2c03b9e551d",
                    "Init": false,
                    "DNSConfig": {},
                    "Isolation": "default"
                },
                "Resources": {
                    "Limits": {},
                    "Reservations": {}
                },
                "Placement": {
                    "Platforms": [
                        {
                            "Architecture": "amd64",
                            "OS": "linux"
                        },
                        {
                            "OS": "linux"
                        },
                        {
                            "Architecture": "arm64",
                            "OS": "linux"
                        }
                    ]
                },
                "ForceUpdate": 0,
                "Runtime": "container"
            },
            "Mode": {
                "Replicated": {
                    "Replicas": 1
                }
            },
            "EndpointSpec": {
                "Mode": "vip",
                "Ports": [
                    {
                        "Protocol": "tcp",
                        "TargetPort": 5000,
                        "PublishedPort": 5000,
                        "PublishMode": "ingress"
                    }
                ]
            }
        },
        "Endpoint": {
            "Spec": {
                "Mode": "vip",
                "Ports": [
                    {
                        "Protocol": "tcp",
                        "TargetPort": 5000,
                        "PublishedPort": 5000,
                        "PublishMode": "ingress"
                    }
                ]
            },
            "Ports": [
                {
                    "Protocol": "tcp",
                    "TargetPort": 5000,
                    "PublishedPort": 5000,
                    "PublishMode": "ingress"
                }
            ],
            "VirtualIPs": [
                {
                    "NetworkID": "yj8219fu3uq1auabkvbx6gtqx",
                    "Addr": "10.0.0.9/24"
                }
            ]
        }
    },
]
*/
