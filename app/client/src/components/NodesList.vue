<template>
  <v-data-table
    :headers="headers"
    :items="nodes"
    hide-default-header
    hide-default-footer
  >
    <template v-slot:[`item.action`]="{ item }">
  
      <v-menu absolute>
        <template v-slot:activator="{ on, attrs }">
          <v-btn v-bind="attrs" v-on="on">
            <v-icon>mdi-dots-horizontal</v-icon>
          </v-btn>
        </template>
  
        <v-list>
          <v-list-item @click="() => item /* TODO */">
            <v-list-item-title>View logs</v-list-item-title>
          </v-list-item>
          <v-list-item @click="() => item /* TODO */">
            <v-list-item-title>Restart</v-list-item-title>
          </v-list-item>
          <v-list-item @click="() => item /* TODO */">
            <v-list-item-title>Disable</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
  
    </template>
  </v-data-table>
</template>

<script>
  export default {
    data () {
      return {
        headers: [
          {
            text: 'Node',
            align: 'start',
            sortable: false,
            value: 'name',
          },
          {
            text: 'Action',
            align: 'end',
            sortable: false,
            value: 'action'
          },
        ],
        nodes: [
          {
            name: 'Raspberry Pi 3B+',
          },
          {
            name: 'Windows 10',
          },
          {
            name: 'Ubuntu 20.04',
          },
        ]
      }
    },
  }
</script>
