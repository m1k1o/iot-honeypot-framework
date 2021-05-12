<template>
  <v-card>
    <v-card-title>
      <span class="headline">Proxy specification</span>
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col>
          <v-select
            label="Service hostname"
            v-model="data.service"
            :items="servicesList"
          ></v-select>
        </v-col>
        <div style="padding-top: 2.5em;">
          <v-btn @click="LoadServices()" icon color="green"><v-icon>mdi-refresh</v-icon></v-btn>
        </div>
      </v-row>
      <v-row>
        <v-col>
          <v-select
            label="Node hostname"
            v-model="data.node"
            :items="nodesList"
          ></v-select>
        </v-col>
        <div style="padding-top: 2.5em;">
          <v-btn @click="LoadNodes()" icon color="green"><v-icon>mdi-refresh</v-icon></v-btn>
        </div>
      </v-row>
    </v-card-text>
    <v-card-title>
      <span class="headline">Ports Specification</span>
    </v-card-title>
    <v-card-text>
      <v-row v-for="port, index in data.ports" :key="index">
        <v-col>
          <v-select
            label="Protocol"
            :value="port.protocol"
            @input="$set(data.ports, index, { ...port, protocol: String($event) })"
            :items="[ 'tcp', 'udp' ]"
          ></v-select>
        </v-col>
        <v-col>
          <v-text-field
            label="Target Port"
            :value="port.targetPort"
            @input="$set(data.ports, index, { ...port, targetPort: Number($event) })"
          ></v-text-field>
        </v-col>
        <div style="padding-top: 2.5em;">
          <v-icon>mdi-arrow-right</v-icon>
        </div>
        <v-col>
          <v-text-field
            label="Published Port"
            :value="port.publishedPort"
            @input="$set(data.ports, index, { ...port, publishedPort: Number($event) })"
          ></v-text-field>
        </v-col>
        <div style="padding-top: 2.5em;">
          <v-btn @click="RemovePort(index)" icon color="red"><v-icon>mdi-close</v-icon></v-btn>
        </div>
      </v-row>
      <v-btn @click="AddPort" icon color="green"><v-icon>mdi-plus</v-icon></v-btn>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        color="blue darken-1"
        text
        @click="Finish"
      >
        Close
      </v-btn>
      <v-btn
        color="blue"
        @click="Create"
      >
        Create
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'

import { ProxySpec } from '@/api/index.ts'

  export default Vue.extend({
    computed: {
      ...mapGetters([
        'nodes',
        'services',
      ]),
      servicesList(): { text: string; value: string } {
        return this.services.map(({ name }: { name: string }) => ({ text: name, value: name }))
      },
      nodesList(): { text: string; value: string } {
        return this.nodes.map(({ hostname, role }: { hostname: string; role: string }) => ({ text: hostname + ' (' + role + ')', value: hostname }))
      },
    },
    data () {
      return {
        data: {} as ProxySpec,
      }
    },
    methods: {
      async Create() {
        try {
          await this.$store.dispatch('PROXY_CREATE', this.data)
          await this.$swal({
            title: 'Proxy created!',
            icon: 'success',
          })

          this.Finish()
        } catch (res) {
          const { message } = await res.json()
          const { value } = await this.$swal({
            title: res.statusText,
            text: message,
            icon: 'error',
            showCancelButton: true,
            confirmButtonText: "Close",
            cancelButtonText: "Edit",
          })
          
          if (value) {
            this.Finish()
          }
        }
      },
      Finish() {
        this.$emit('finished')
        this.Clear()
      },
      Clear() {
        this.data = {
          service: '',
          node: '',
          ports: [],
        }
      },
      AddPort() {
        this.data.ports?.push({
          protocol: undefined,
          targetPort: undefined,
          publishedPort: undefined,
        })
      },
      RemovePort(index: number) {
        this.data.ports?.splice(index, 1);
      },
      LoadServices() {
        this.$store.dispatch('SERVICES_LOAD')
      },
      LoadNodes() {
        this.$store.dispatch('NODES_LOAD')
      },
    },
    beforeMount() {
      this.Clear()
    }
  })
</script>
