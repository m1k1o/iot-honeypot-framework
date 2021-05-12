<template>
  <v-card>
    <v-card-title>
      <span class="headline">Service specification</span>
    </v-card-title>
    <v-card-text>
      <v-text-field
        label="Name"
        v-model="data.name"
      ></v-text-field>
      <v-text-field
        label="Replicas"
        type="number"
        :value="data.replicas"
        @input="$set(data, 'replicas', Number($event))"
      ></v-text-field>
    </v-card-text>
    <v-card-title>
      <span class="headline">Container Specification</span>
    </v-card-title>
    <v-card-text>
      <v-text-field
        label="Image"
        :value="data.containerSpec.image"
        @input="$set(data.containerSpec, 'image', $event)"
      ></v-text-field>

      <a @click="extended = !extended"> {{ extended ? 'Hide' : 'View' }} extended settings </a>
      <template v-if="extended">
        <v-textarea
          label="Command"
          :value="data.containerSpec.cmd.join('\n')"
          @input="$set(data.containerSpec, 'cmd', $event == '' ? [] : $event.split('\n'))"
          hint="Separated by new line"
        ></v-textarea>
        <v-textarea
          label="Arguments"
          :value="data.containerSpec.args.join('\n')"
          @input="$set(data.containerSpec, 'args', $event == '' ? [] : $event.split('\n'))"
          hint="Separated by new line"
        ></v-textarea>
        <v-text-field
          label="Hostname"
          :value="data.containerSpec.hostname"
          @input="$set(data.containerSpec, 'hostname', $event)"
        ></v-text-field>
        <v-textarea
          label="Environment variables"
          :value="data.containerSpec.env.join('\n')"
          @input="$set(data.containerSpec, 'env', $event == '' ? [] : $event.split('\n'))"
          hint="Separated by new line"
        ></v-textarea>
        <v-text-field
          label="Directory"
          :value="data.containerSpec.dir"
          @input="$set(data.containerSpec, 'dir', $event)"
        ></v-text-field>
      </template>
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

import {
  ServiceSpec,
} from '@/api/index.ts'

  export default Vue.extend({
    data () {
      return {
        extended: false,
        data: {} as ServiceSpec,
      }
    },
    watch: {
      ['data.containerSpec']() {
        console.log("changed")
      }
    },
    methods: {
      Val<T>(val: T): T {
        //if (T instanceof String) {  
        //  console.log("Yey Striongf")
        //}

        return val
      },
      async Create() {
        try {
          await this.$store.dispatch('SERVICE_CREATE', this.data)
          await this.$swal({
            title: 'Service created!',
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
          name: "My_First_Container",
            containerSpec: {
              image: "nginx:latest",
              cmd: [],
              args: [],
              hostname: "",
              env: [],
              dir: "",
            },
          replicas: 1
        }
      },
    },
    beforeMount() {
      this.Clear()
    }
  })
</script>
