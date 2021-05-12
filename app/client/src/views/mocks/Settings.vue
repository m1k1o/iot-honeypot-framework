<template>
  <v-container>
    <v-sparkline
      :value="s_value"
      type="bar"
      color="rgba(255, 255, 255, .7)"
      height="25"
    >
    </v-sparkline>

    <v-row no-gutters>
      <v-col>
        <v-card
          class="pa-2 ma-2"
          outlined
          tile
        >
          <v-toolbar-title class="pa-2">Connection settings</v-toolbar-title>
          <v-card-text>
            <v-text-field
              label="Public IP address"
              value=""
            ></v-text-field>
            <v-text-field
              label="SSH Port"
              value=""
            ></v-text-field>
            <v-text-field
              label="SSH Username"
              value=""
            ></v-text-field>
            <v-text-field
              label="SSH Password"
              value=""
            ></v-text-field>
            <v-btn>Generate SSH key</v-btn>
            <v-checkbox
              label="Allow password login"
              :value="true"
            ></v-checkbox>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col>
        <v-card
          class="pa-2 ma-2"
          outlined
          tile
        >
          <v-toolbar-title class="pa-2">Monitoring</v-toolbar-title>
          <v-card-text>
            <v-text-field
              label="Metrics port"
              value="9001"
            ></v-text-field>
            <v-text-field
              label="Metrics url"
              value="/_metrics"
            ></v-text-field>
            <v-btn>Downlaod grafana Dashboard</v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
