<template>
  <v-container>
    <v-data-table
      :headers="headers"
      :items="services"
      class="elevation-1"
    >
      <template v-slot:[`item.name`]="{ item }">
        {{ item.name }}
      </template>
      <template v-slot:[`item.container`]="{ item }">
        {{ item.containerSpec.image }}
      </template>
      <template v-slot:[`item.status`]="{ item }">
        <v-chip v-if="item.status.running == 0" color="red" dark> {{ item.status.running }} / {{ item.status.desired }} </v-chip>
        <v-chip v-else-if="item.status.running != item.status.desired" color="yellow" dark> {{ item.status.running }} / {{ item.status.desired }} </v-chip>
        <v-chip v-else color="green" dark> {{ item.status.running }} / {{ item.status.desired }} </v-chip>
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
      <ServiceCreate
        @finished="dialog = false"
      />
    </v-dialog>
    
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import ServiceCreate from '@/components/ServiceCreate.vue'

  export default Vue.extend({
    computed: {
      ...mapGetters([
        'services',
      ]),
    },
    data () {
      return {
        dialog: null,
        headers: [
          {
            text: 'Name',
            sortable: true,
            value: 'name'
          },
          {
            text: 'Container',
            sortable: false,
            value: 'container',
          },
          {
            text: 'Replicas',
            sortable: false,
            value: 'status'
          },
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
        this.$store.dispatch('SERVICES_LOAD')
      },
      async Remove(serviceId: string) {
        const { value } = await this.$swal({
          title: "Remove service",
          text: "Do you really want to remove this service?",
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: "Yes",
          cancelButtonText: "No",
        })
        
        if (value) {
          try {
            await this.$store.dispatch('SERVICE_REMOVE', serviceId)
            await this.$swal({
              title: 'Service removed!',
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
      ServiceCreate
    },
  })
</script>
