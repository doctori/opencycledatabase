<template>
  <div id="app">
    <p><img src="./assets/logo.png"></p>
    <router-link to="/">Home</router-link>
    <router-link to="/standards">Standards</router-link>
    <router-view></router-view>
    <footer>
    
  </footer>
  </div>
  
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
