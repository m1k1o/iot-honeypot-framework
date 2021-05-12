<template>
  <v-container>
    <v-data-table
      :headers="headers"
      :items="proxies"
      class="elevation-1"
      show-group-by
    >
      <template v-slot:[`item.service`]="{ item }">
        {{ item.service }}
      </template>
      <template v-slot:[`item.node`]="{ item }">
        {{ item.node }}
      </template>
      <template v-slot:[`item.ports`]="{ item }">
        <v-chip
          color="blue"
          dark class="ma-1"
          v-for="{ protocol, targetPort, publishedPort } in item.ports"
          :key="protocol + targetPort + publishedPort">
          {{ targetPort }}
          <v-icon>mdi-arrow-right</v-icon>
          {{ publishedPort }}/{{ protocol }}
        </v-chip>
      </template>
      <template v-slot:[`item.status`]="{ item }">
        <v-chip v-if="item.running" color="green" dark> up </v-chip>
        <v-chip v-else color="red" dark> down </v-chip>
      </template>
      <template v-slot:[`item.action`]="{ item }">

        <v-menu absolute>
          <template v-slot:activator="{ on, attrs }">
            <v-btn v-bind="attrs" v-on="on">
              <v-icon>mdi-dots-horizontal</v-icon>
            </v-btn>
          </template>

          <v-list>
            <!--
            <v-list-item @click="/* TODO */">
              <v-list-item-title>Detail</v-list-item-title>
            </v-list-item>
            <v-list-item @click="/* TODO */">
              <v-list-item-title>Stats</v-list-item-title>
            </v-list-item>
            <v-divider />
            -->
            <v-list-item @click="Remove(item.id)">
              <v-list-item-title>Remove</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>

      </template>
    </v-data-table>
    <v-btn @click="Load" class="mt-3" color="green" icon style="position: absolute;margin-top: -46px!important;margin-left: 10px;"><v-icon>mdi-refresh</v-icon></v-btn>

    <v-dialog
      v-model="dialog"
      persistent
      max-width="600px"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-btn 
          v-bind="attrs"
          v-on="on"
          class="mt-3"
          color="success"
        > + Add new </v-btn>
      </template>
      <ProxyCreate
        @finished="dialog = false"
      />
    </v-dialog>

  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import ProxyCreate from '@/components/ProxyCreate.vue'

  export default Vue.extend({
    computed: {
      ...mapGetters([
        'proxies',
      ]),
    },
    data () {
      return {
        dialog: null,
        headers: [
          {
            text: 'Service',
            sortable: true,
            groupable: true,
            value: 'service'
          },
          {
            text: 'Node',
            sortable: true,
            groupable: true,
            value: 'node',
          },
          {
            text: 'Ports',
            sortable: false,
            groupable: false,
            value: 'ports'
          },
          {
            text: 'Status',
            sortable: false,
            groupable: false,
            value: 'status'
          },
          {
            text: 'Action',
            align: 'end',
            sortable: false,
            groupable: false,
            value: 'action'
          },
        ],
      }
    },
    methods: {
      Load() {
        this.$store.dispatch('PROXIES_LOAD')
      },
      async Remove(proxyId: string) {
        const { value } = await this.$swal({
          title: "Remove proxy",
          text: "Do you really want to remove this proxy?",
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: "Yes",
          cancelButtonText: "No",
        })
        
        if (value) {
          try {
            await this.$store.dispatch('PROXY_REMOVE', proxyId)
            await this.$swal({
              title: 'Proxy removed!',
              icon: 'success',
            })
          } catch (res) {
            const { message } = await res.json()
            await this.$swal({
              title: res.statusText,
              text: message,
              icon: 'error',
            })
          }
        }
      }
    },
    mounted() {
      this.Load()
    },
    components: {
      ProxyCreate
    },
  })
</script>
