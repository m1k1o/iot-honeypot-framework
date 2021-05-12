<template>
  <v-container>
    <v-card
      class="pa-2 ma-2"
      outlined
      tile
    >
      <v-toolbar-title class="pa-2">Add new worker</v-toolbar-title>
      <v-card-text>
        <div class="pb-4">
          <h2 class="pb-2">Step 1: Install dependencies</h2>
          <p>OpenVPN and Docker need to installed on target server.</p>
          <code class="pa-3">{{ [
            'sudo apt update',
            'sudo apt install openvpn docker.io'
          ].join('\n') }}</code>
        </div>

        <div class="pb-4">
          <h2 class="pb-2">Step 2: Connect to VPN</h2>
          <p>Node needs to be connected to OpenVPN, so it can reach master node.</p>
          <code class="pa-3" v-if="ovpnClientConf">{{ [
            'echo \'' + ovpnClientConf.replace('\'', '\\\'') + '\' > /etc/openvpn/client.conf',
            'sudo systemctl start openvpn@client'
          ].join('\n') }}</code>
          <!--<v-btn v-else @click="GenerateOpenVPN" color="info"> Generate code </v-btn>-->
        </div>

        <div class="pb-4">
          <h2 class="pb-2">Step 3: Join docker swarm</h2>
          <p>Node needs to join docker swarm as a worker.</p>
          <code class="pa-3 mb-3" v-if="swarmToken && swarmURL">{{ [
            'docker swarm join --listen-addr tun0 --token ' + swarmToken + ' ' + swarmURL
          ].join('\n') }}</code>
          <v-btn @click="GenerateSwarmInvitation" color="info"> Recreate invitation </v-btn>
        </div>

        <div>
          <h2 class="pb-2">Done</h2>
          <p class="mb-0">Your new worker should be connected now.</p>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'

  export default Vue.extend({
    data () {
      return {
        ovpnClientConf: '',
        swarmToken: '',
        swarmURL: '',
      }
    },
    methods: {
      GenerateOpenVPN() {
        this.ovpnClientConf = `client
dev tun
proto udp
remote ${location.hostname} 1194
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
auth SHA512
cipher AES-256-CBC
ignore-unknown-option block-outside-dns
block-outside-dns
verb 3
pull-filter ignore "redirect-gateway"
route 10.8.0.0 255.255.255.0
push "route 10.8.0.0 255.255.255.0"
<ca>
-----BEGIN CERTIFICATE-----
foobarfoobarfoobarfoobarfoobarfoobar
-----END CERTIFICATE-----
</ca>
<cert>
-----BEGIN PRIVATE KEY-----
foobarfoobarfoobarfoobarfoobarfoobar
-----END PRIVATE KEY-----
</key>
<tls-crypt>
-----BEGIN OpenVPN Static key V1-----
foobarfoobarfoobarfoobarfoobarfoobar
-----END OpenVPN Static key V1-----
</tls-crypt>`
      },
      async GenerateSwarmInvitation() {
        const { token, addr } = await this.$store.dispatch('NODE_JOIN')
        this.swarmToken = token
        this.swarmURL = addr
      }
    },
    mounted() {
      this.GenerateSwarmInvitation()
    }
  })
</script>

<style scoped>
  code {
    display:block;
    white-space:break-spaces;
    font-size:120%;
    font-family:Consolas,monospace;
  }
</style>
