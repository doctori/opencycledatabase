<template>
  <v-container>
    <v-container v-if="loading" class="loading">
        Loading...
    </v-container>
    <v-container v-if="error" class="error">
          {{ error.message }}
    </v-container>
  <v-container>
    <h4>Welcome to the Open Cycle database</h4>
    You should find here all the component and compatibility matrix, 
    that should help you know what components are compatible with your bike. (soon)
  </v-container>
  </v-container>
</template>


<script>

export default {
  name: 'Main',
  components: {
  },
  data(){
    return{

      loading: false,
      post: null,
      error: null,
      createStandard: false,
      createBrand: false,
      standards: []
    }
  },
  // when the template is created
  created(){
    this.fetchData()
  },
  watch: {
    // Call again the method if the route changes
    '$route':'fetchData'
  },
  methods: {
    enableCreateStandard(){
      this.createStandard = !this.createStandard;
      this.createBrand = false
    },
    fetchData(){
      this.error = this.post = null;
      this.loading = true;
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
