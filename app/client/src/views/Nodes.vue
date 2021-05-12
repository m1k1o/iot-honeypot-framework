<template>
  <v-container>
    <v-data-table
      :headers="headers"
      :items="nodes"
      class="elevation-1"
    >
      <template v-slot:[`item.role`]="{ item }">
        <v-chip color="green" dark v-if="item.role == 'manager'">manager</v-chip>
        <v-chip color="gray" dark v-else>{{ item.role }}</v-chip>
      </template>
      <template v-slot:[`item.hostname`]="{ item }">
        {{ item.hostname }}
      </template>
      <template v-slot:[`item.platform`]="{ item }">
        {{ item.platform.architecture }} / {{ item.platform.os }}
      </template>
      <template v-slot:[`item.resources`]="{ item }">
      {{ Math.round(item.resources.nanoCpus / 10e8) / 10 }} /
      {{ Math.round(item.resources.memoryBytes / 10e8) / 10 }}
      </template>
      <template v-slot:[`item.status`]="{ item }">
        <v-chip color="green" dark v-if="item.status.state == 'ready'">ready</v-chip>
        <v-chip color="red" dark v-else-if="item.status.state == 'down'">down</v-chip>
        <v-chip color="gray" dark v-else>{{ item.status.state }}</v-chip>
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

    <v-btn :to="{ name: 'node-join' }" class="mt-3" color="success"> + Add new </v-btn>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'

  export default Vue.extend({
    computed: {
      ...mapGetters([
        'nodes',
      ]),
    },
    data () {
      return {
        headers: [
          { text: 'Hostname', value: 'hostname' },
          { text: 'Role', value: 'role' },
          { text: 'Platform', value: 'platform' },
          { text: 'Resources', value: 'resources' },
          { text: 'Status', value: 'status' },
          {
            text: 'Action',
            align: 'end',
            sortable: false,
            value: 'action'
          },
        ],
      }
    },
    methods: {
      Load() {
        this.$store.dispatch('NODES_LOAD')
      },
      async Remove(nodeId: string) {
        const { value } = await this.$swal({
          title: "Remove node",
          text: "Do you really want to remove this node?",
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: "Yes",
          cancelButtonText: "No",
        })
        
        if (value) {
          try {
            await this.$store.dispatch('NODE_REMOVE', nodeId)
            await this.$swal({
              title: 'Node removed!',
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
  })
</script>
