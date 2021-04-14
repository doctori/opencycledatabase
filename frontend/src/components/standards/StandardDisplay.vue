<template>
  <v-expansion-panels>
    <v-expansion-panel>
      <v-expansion-panel-header>
        {{$t('std.list-name')}}
      </v-expansion-panel-header>
      <v-expansion-panel-content>
    <v-row class="standard-display">
      <v-col >
      <ul>
        <li v-for="standard in standards" v-bind:key="standard.Type" v-bind:title="standard.Type">{{standard.Type}}</li>
      </ul>
      </v-col>
    </v-row>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
export default {
  name: 'StandardDisplay',
  props: {'standards': Array }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
a {
  color: #42b983;
}
</style>
