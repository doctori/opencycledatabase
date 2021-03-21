<template>
<div class="standard-display">
  <ul>
    <li v-for="standard in standards" v-bind:key="standard.Type" v-bind:title="standard.Type">{{standard.Type}}</li>
  </ul>
</div>
</template>

<script>
export default {
  name: 'StandardDisplay',
  props: {'standards':Array}
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
