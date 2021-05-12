<template>
  <v-data-table
    :headers="_headers"
    :items="containers"
    hide-default-footer
  >
    <template v-slot:[`item.ports`]="{ item }">
      <v-chip color="blue" dark v-for="port in item.ports" :key="port">{{ port }}</v-chip>
    </template>
    <template v-slot:[`item.started`]="{ item }">
      {{ item.started | timeago }}
    </template>

    <template v-slot:[`item.status`]="{ item }">
      <v-tooltip bottom v-if="item.status === 'stopped'">
        <template v-slot:activator="{ on, attrs }">
          <v-icon color="grey" style="cursor: auto;" v-bind="attrs" v-on="on">mdi-pause-circle</v-icon>
        </template>
        <span>Stopped</span>
      </v-tooltip>
      <v-tooltip bottom v-if="item.status === 'failed'">
        <template v-slot:activator="{ on, attrs }">
          <v-icon color="red" style="cursor: auto;" v-bind="attrs" v-on="on">mdi-alert-circle</v-icon>
        </template>
        <span>Failed</span>
      </v-tooltip>
      <v-tooltip bottom v-if="item.status === 'running'">
        <template v-slot:activator="{ on, attrs }">
          <v-icon color="green" style="cursor: auto;" v-bind="attrs" v-on="on">mdi-play-circle</v-icon>
        </template>
        <span>Running</span>
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
          <template v-if="item.status === 'running'">
            <v-list-item @click="() => item /* TODO */">
              <v-list-item-title>Logs</v-list-item-title>
            </v-list-item>
            <v-list-item @click="() => item /* TODO */">
              <v-list-item-title>Execute</v-list-item-title>
            </v-list-item>
            <v-divider />
            <v-list-item @click="() => item /* TODO */">
              <v-list-item-title>Restart</v-list-item-title>
            </v-list-item>
            <v-list-item @click="() => item /* TODO */">
              <v-list-item-title>Stop</v-list-item-title>
            </v-list-item>
          </template>
          <template v-else>
            <v-list-item @click="() => item /* TODO */">
              <v-list-item-title>Start</v-list-item-title>
            </v-list-item>
            <v-list-item @click="() => item /* TODO */">
              <v-list-item-title>Remove</v-list-item-title>
            </v-list-item>
          </template>
        </v-list>
      </v-menu>
  
    </template>
  </v-data-table>
</template>

<script>
import moment from 'moment'

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
            value: 'container',
          },
          {
            text: 'Node',
            value: 'node',
          },
          { text: 'Started', value: 'started' },
          {
            text: 'Exposed ports',
            align: 'start',
            sortable: false,
            value: 'ports'
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
        containers: [
          {
            container: 'SSH Honeypot',
            node: 'Raspberry Pi 3B+',
            ports: [ '22/tcp' ],
            started: moment().subtract(28, 'seconds'),
            status: 'running',
          },
          {
            container: 'Router Honeypot',
            node: 'Windows 10',
            ports: [ '8080/tcp' ],
            started: moment().subtract(10, 'hours'),
            status: 'failed',
          },
          {
            container: 'Open WebIf',
            node: 'Ubuntu 20.04',
            ports: [ '8888/tcp' ],
            started: moment().subtract(4, 'minutes'),
            status: 'stopped',
          },
          {
            container: 'IP Camera',
            node: 'Raspberry Pi 3B+',
            ports: [ '81/tcp' ],
            started: moment().subtract(30, 'minutes'),
            status: 'stopped',
          },
          {
            container: 'HTTP Server',
            node: 'Windows 10',
            ports: [ '80/tcp' ],
            started: moment().subtract(2, 'hours'),
            status: 'running',
          },
          {
            container: 'HTTP Server',
            node: 'Ubuntu 20.04',
            ports: [ '80/tcp' ],
            started: moment().subtract(2, 'hours'),
            status: 'running',
          },
        ]
      }
    },
  }
</script>
