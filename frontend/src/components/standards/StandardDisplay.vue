<template>
  <v-row class="standard-display">
    <v-col cols="3">
    {{ standardInput.Name }}
    </v-col>
    <v-col >
    {{ standardInput.Description }}
    </v-col>
  </v-row>
</template>

<script>

export default {
  name: 'StandardDisplay',
  props: {'standardInput': Object }
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
