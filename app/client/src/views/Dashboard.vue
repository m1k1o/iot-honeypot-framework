<template>
  <v-container>
    <v-row>
      <v-col lg="3" cols="sm" class="pb-2">
        <v-card to="/services">
          <v-row class="no-gutters">
            <div class="col-auto">
              <div class="green fill-height">&nbsp;</div>
            </div>
            <div class="col pa-3 py-4 green--text">
              <h5 class="text-truncate text-uppercase">Services / Replicas</h5>
              <h1>{{ services.length }}</h1>
            </div>
          </v-row>
        </v-card>
      </v-col>
      <v-col lg="3" cols="sm" class="pb-2">
        <v-card to="/services">
          <v-row class="no-gutters">
            <div class="col-auto">
              <div class="blue fill-height">&nbsp;</div>
            </div>
            <div class="col pa-3 py-4 blue--text">
              <h5 class="text-truncate text-uppercase">Replicas <small>running / desired</small></h5>
              <h1>{{ replicasRunning }} / {{ replicasDesired }}</h1>
            </div>
          </v-row>
        </v-card>
      </v-col>
      <v-col lg="3" cols="sm" class="pb-2">
        <v-card to="/proxies">
          <v-row class="no-gutters">
            <div class="col-auto">
              <div class="purple fill-height">&nbsp;</div>
            </div>
            <div class="col pa-3 py-4 cyan--text">
              <h5 class="text-truncate text-uppercase">Proxies</h5>
              <h1>{{ proxies.length }}</h1>
            </div>
          </v-row>
        </v-card>
      </v-col>
      <v-col lg="3" cols="sm" class="pb-2">
        <v-card to="/nodes">
          <v-row class="no-gutters">
            <div class="col-auto">
              <div class="cyan fill-height">&nbsp;</div>
            </div>
            <div class="col pa-3 py-4 purple--text">
              <h5 class="text-truncate text-uppercase">Nodes</h5>
              <h1>{{ nodes.length }}</h1>
            </div>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col 
        v-for="{ name, icon } in items"
        :key="name"
        lg="4"
        cols="sm"
        class="pb-2"
      >
        <a :href="ServiceUrl(name)" style="text-decoration:none;" target="_blank">
          <v-card @click="() => {}">
            <v-row class="no-gutters">
              <div class="col-auto d-flex align-center ma-4 mr-0">
                <v-icon large class="teal--text text--lighten-1">{{ icon }}</v-icon>
              </div>
              <div class="col ma-4 text-uppercase">
                <h2>{{ name }}</h2>
                <h5 class="teal--text text--lighten-1">{{ name }}.{{ host }}</h5>
              </div>
            </v-row>
          </v-card>
        </a>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'

  export default Vue.extend({
    computed: {
      ...mapGetters([
        'services',
        'proxies',
        'nodes',
      ]),
      replicasRunning() {
        return this.services.reduce((total: number, { status: { running } }: { status: { running: number } }) => total + running, 0)
      },
      replicasDesired() {
        return this.services.reduce((total: number, { status: { desired } }: { status: { desired: number } }) => total + desired, 0)
      },
      host() {
        return location.host
      }
    },
    data () {
      return {
        items: [
          {
            icon: 'mdi-database',
            name: 'kibana',
          },
          {
            icon: 'mdi-docker',
            name: 'portainer',
          },
          {
            icon: 'mdi-thermostat-box',
            name: 'grafana',
          },
          {
            icon: 'mdi-database-refresh-outline',
            name: 'prometheus',
          },
          {
            icon: 'mdi-alert',
            name: 'alertmanager',
          },
          {
            icon: 'mdi-lock-outline',
            name: 'traefik',
          },
        ],
      }
    },
    methods: {
      ServiceUrl(name: string): string {
        return location.protocol + '//' + name + '.' + location.host
      }
    },
    mounted() {
      this.$store.dispatch('SERVICES_LOAD')
      this.$store.dispatch('PROXIES_LOAD')
      this.$store.dispatch('NODES_LOAD')
    }
  })
</script>
