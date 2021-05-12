<template>
  <v-data-table
    :headers="_headers"
    :items="nodes"
    hide-default-footer
  >
    <template v-slot:[`item.ports`]="{ item }">
      <v-chip color="blue" dark>{{ item.containerPort }}</v-chip>
      <v-icon>mdi-arrow-right</v-icon>
      <v-chip color="blue" dark>{{ item.nodePort }}</v-chip>
    </template>

    <template v-slot:[`item.status`]="{ item }">
      <v-tooltip bottom v-if="item.enabled">
        <template v-slot:activator="{ on, attrs }">
          <v-icon color="green" style="cursor: auto;" v-bind="attrs" v-on="on">mdi-checkbox-marked-circle</v-icon>
        </template>
        <span>Enabled</span>
      </v-tooltip>
      <v-tooltip bottom v-else>
        <template v-slot:activator="{ on, attrs }">
          <v-icon color="grey" style="cursor: auto;" v-bind="attrs" v-on="on">mdi-close-circle</v-icon>
        </template>
        <span>Disabled</span>
      </v-tooltip>
    </template>

    <template v-slot:[`item.action`]="{ item }">
  
      <v-menu absolute>
        <template v-slot:activator="{ on, attrs }">
          <v-btn v-bind="attrs" v-on="on">
            <v-icon>mdi-dots-horizontal</v-icon>
          </v-btn>
        </template>
  
        <v-list>
          <v-list-item @click="() => item /* TODO */">
            <v-list-item-title>Edit</v-list-item-title>
          </v-list-item>
          <v-list-item @click="() => item /* TODO */" v-if="item.enabled">
            <v-list-item-title>Disable</v-list-item-title>
          </v-list-item>
          <v-list-item @click="() => item /* TODO */" v-else>
            <v-list-item-title>Enable</v-list-item-title>
          </v-list-item>
          <v-list-item @click="() => item /* TODO */">
            <v-list-item-title>Remove</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </template>
  </v-data-table>
</template>

<script>
  export default {
    props: {
      'no-container': {
        type: Boolean,
        default: false
      },
      'no-node': {
        type: Boolean,
        default: false
      }
    },
    computed: {
      _headers() {
        let headers = [...this.headers];

        if(this.$props.noContainer) {
          headers = headers.filter(({ value }) => value != 'container');
        }

        if(this.$props.noNode) {
          headers = headers.filter(({ value }) => value != 'node');
        }

        return headers;
      }
    },
    data () {
      return {
        headers: [
          {
            text: 'Container',
            align: 'start',
            sortable: false,
            value: 'container',
          },
          {
            text: 'Published ports',
            align: 'start',
            sortable: false,
            value: 'ports',
          },
          {
            text: 'Node',
            align: 'start',
            sortable: false,
            value: 'node'
          },
          {
            text: 'Status',
            align: 'center',
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
        nodes: [
          {
            container: 'SSH Honeypot',
            containerPort: '80/tcp',
            nodePort: '8080/tcp',
            node: 'Raspberry Pi 3B+',
            enabled: true,
          },
          {
            container: 'SSH Honeypot',
            containerPort: '80/tcp',
            nodePort: '8080/tcp',
            node: 'Windows 10',
            enabled: false,
          },
          {
            container: 'SSH Honeypot',
            containerPort: '80/tcp',
            nodePort: '8080/tcp',
            node: 'Ubuntu 20.04',
            enabled: true,
          },
        ]
      }
    },
  }
</script>
