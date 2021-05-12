# IoT Honeypot Framework
Scalable Honeypot Framework using Docker Swarm and Elasticsearch.

# Installation

## 1. Set up VPN.
Create OVPN server on master and connect all workers to this VPN network.

### On master
Go through installation wizard:

```sh
./openvpn.sh
```

### On worker(s)
Add new client on master and copy its configuration to worker.

```sh
# Install openvpn client
sudo apt install openvpn

# Copy generated configuration from master
mv ovpn-client.conf /etc/openvpn/client.conf

# Start client configuration
sudo systemctl start openvpn@client
```

## 2. Set up Docker swarm.
On all devices install docker first:

```sh
sudo apt install docker.io
sudo systemctl enable docker.service
```

### On master
Master will be swarm manager. Force it to use newly created VPN:

```sh
docker swarm init --listen-addr tun0 --advertise-addr 10.8.0.1
```

### On worker(s)
Worker will join to swarm cluster created by master. Do not forget, to force worker as well to use newly created VPN:

```sh
docker swarm join --listen-addr tun0 --token <created-token> 10.8.0.1:2377
```

## 3. Set up Docker registry.
Registry needs to be accessible from all nodes. Used IP `10.8.0.1` is `tun0` IP address of master node.

### On master
This script will generate new certificates and run docker registry.

```sh
./docker-registry/generate_certs.sh
sudo ./docker-registry/create.sh "10.8.0.1"
```

#### Regenerate certificates
Certificates can be easily replaced on all nodes by runnin following commands.

```sh
./docker-registry/generate_certs.sh
sudo ./docker-registry/swarm_update.sh "10.8.0.1"
```

### On worker(s)
You need to copy `certs` folder to all workers along with `apply.sh`.

```sh
sudo ./apply.sh "10.8.0.1"
```

## 4. Publish Docker proxy.
Docker proxy needs to be built and published into docker registry.

```sh
./docker-proxy/build.sh
./docker-proxy/publish.sh
```

## 5. Publish Manager App.
Manager App needs to be built and published into docker registry.

```sh
./app/build.sh
./app/publish.sh
```


## 6. Deploy master stack
Before installation, `.env` file must be properly filled with all relevant data.

```sh
cd master-stack
cp .env.example .env
```

Then, master stack can  be installed.

```sh
cd master-stack
./install.sh
```

And deployed:

```sh
./deploy.sh base
./deploy.sh logs
./deploy.sh app
./deploy.sh portainer
./deploy.sh grafana
```

You will need this for elasticsearch:

```sh
echo 'vm.max_map_count = 262144' >> /etc/sysctl.conf
sysctl -w vm.max_map_count=262144
systemctl restart docker
```
