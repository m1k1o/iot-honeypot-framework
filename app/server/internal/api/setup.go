package api

// set up docker registry
// docker service create --name registry --publish published=5000,target=5000 --constraint node.role==manager registry:2

// create my attachable overlay
// docker network create -d overlay --attachable my-attachable-overlay

// build docker proxy
// docker build -t localhost:5000/docker_proxy "${REPO_ROOT}/docker-proxy"
// docker push localhost:5000/docker_proxy
