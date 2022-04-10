package api

import (
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/go-chi/chi"

	"m1k1o/ioth/internal/utils"
)

type NodeStatus struct {
	Addr  string          `json:"addr,omitempty"`
	State swarm.NodeState `json:"state,omitempty"`
}

type Resources struct {
	NanoCPUs    int64 `json:"nano_cpus,omitempty"`
	MemoryBytes int64 `json:"memory_bytes,omitempty"`
}

type Platform struct {
	Architecture string `json:"architecture,omitempty"`
	OS           string `json:"os,omitempty"`
}

type NodeSpec struct {
	ID        string         `json:"id,omitempty"`
	Role      swarm.NodeRole `json:"role,omitempty"`
	Hostname  string         `json:"hostname,omitempty"`
	Platform  Platform       `json:"platform,omitempty"`
	Resources Resources      `json:"resources,omitempty"`
	Status    NodeStatus     `json:"status,omitempty"`
}

type NodeJoin struct {
	Token string `json:"token,omitempty"`
	Addr  string `json:"addr,omitempty"`
}

func (a *ApiManagerCtx) nodes(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		swarmNodes, err := a.cli.NodeList(r.Context(), types.NodeListOptions{})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		nodes := []NodeSpec{}
		for _, node := range swarmNodes {
			nodes = append(nodes, NodeSpec{
				ID:       node.ID,
				Role:     node.Spec.Role,
				Hostname: node.Description.Hostname,
				Platform: Platform{
					Architecture: node.Description.Platform.Architecture,
					OS:           node.Description.Platform.OS,
				},
				Resources: Resources{
					NanoCPUs:    node.Description.Resources.NanoCPUs,
					MemoryBytes: node.Description.Resources.MemoryBytes,
				},
				Status: NodeStatus{
					State: node.Status.State,
					Addr:  node.Status.Addr,
				},
			})
		}

		utils.HttpSuccess(w, nodes)
	})

	const rotateTokens = true
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		if rotateTokens {
			sw, err := a.cli.SwarmInspect(r.Context())
			if err != nil {
				utils.HttpInternalServer(w, err)
				return
			}

			if err := a.cli.SwarmUpdate(r.Context(), sw.Version, sw.Spec, swarm.UpdateFlags{
				RotateWorkerToken: true,
			}); err != nil {
				utils.HttpInternalServer(w, err)
				return
			}
		}

		// get token
		sw, err := a.cli.SwarmInspect(r.Context())
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		// get node id
		info, err := a.cli.Info(r.Context())
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		// get manager addr
		node, _, err := a.cli.NodeInspectWithRaw(r.Context(), info.Swarm.NodeID)
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		utils.HttpSuccess(w, NodeJoin{
			Token: sw.JoinTokens.Worker,
			Addr:  node.ManagerStatus.Addr,
		})
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := a.cli.NodeRemove(r.Context(), id, types.NodeRemoveOptions{Force: true})
		if err != nil {
			utils.HttpInternalServer(w, err)
			return
		}

		utils.HttpSuccess(w)
	})
}

/*

Sample cli.NodeList Response:

[{
	"ID": "a4zkgeahtud74bdsutisgz6my",
	"Version": {
		"Index": 55
	},
	"CreatedAt": "2020-10-28T11:54:36.922024263Z",
	"UpdatedAt": "2020-10-28T15:40:35.970056628Z",
	"Spec": {
		"Labels": {},
		"Role": "worker",
		"Availability": "active"
	},
	"Description": {
		"Hostname": "iothp-worker",
		"Platform": {
			"Architecture": "x86_64",
			"OS": "linux"
		},
		"Resources": {
			"NanoCPUs": 4000000000,
			"MemoryBytes": 4135698432
		},
		"Engine": {
			"EngineVersion": "19.03.6",
			"Plugins": [ ... ]
		},
		"TLSInfo": { ... }
	},
	"Status": {
		"State": "ready",
		"Addr": "10.8.0.2"
	}
}, {
	"ID": "f4l1h1iv8f6b079f1o2qn2ecx",
	"Version": {
		"Index": 49
	},
	"CreatedAt": "2020-10-28T11:54:21.66840247Z",
	"UpdatedAt": "2020-10-28T15:40:33.70446894Z",
	"Spec": {
		"Labels": {},
		"Role": "manager",
		"Availability": "active"
	},
	"Description": {
		"Hostname": "iothp-manager",
		"Platform": {
			"Architecture": "x86_64",
			"OS": "linux"
		},
		"Resources": {
			"NanoCPUs": 4000000000,
			"MemoryBytes": 4135706624
		},
		"Engine": {
			"EngineVersion": "19.03.6",
			"Plugins": [ ... ]
		},
		"TLSInfo": { ... }
	},
	"Status": {
		"State": "ready",
		"Addr": "10.8.0.1"
	},
	"ManagerStatus": {
		"Leader": true,
		"Reachability": "reachable",
		"Addr": "10.8.0.1:2377"
	}
}]
*/
