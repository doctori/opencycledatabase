<template>
    <header>
      <div class="header">
        <img alt="Vue logo" src="../assets/logo.png">
        <h1>Welcome to the Open Cycle Database</h1>
        
      </div>
    </header>
    <main>
      <div v-if="loading" class="loading">
        Loading...
      </div>
      <div v-if="error" class="error">
        {{ error.message }}
      </div>
      <div v-if="standards" class="content">
        <standard-display :standards="standards"/>
      </div>
      <button id="createStandard" v-on:click="enableCreateStandard()">
        Create Standard
      </button>
      <div v-if="createStandard" class="create-standard">
        <standard-create :standards="standards"/>
      </div>
   </main>

    <footer>
      This is the footer
    </footer>
</template>


<script>


//import Header from './Header.vue'
import StandardDisplay from './StandardDisplay.vue';
import axios from 'axios'
import StandardCreate from './StandardCreate.vue';
export default {
  name: 'Main',
  components: {
    'standard-display':StandardDisplay,
    'standard-create': StandardCreate
  },
  data(){
    return{
      loading: false,
      post: null,
      error: null,
      createStandard: false
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
    this.createStandard = !this.createStandard
  },
  fetchData(){
    this.error = this.post = null
    this.loading = true

    axios
      .get("http://localhost:8080/components")
      .then(response => (this.components = response.data));
    axios
      .get("http://localhost:8080/standards")
      .then(response => (this.standards = response.data ))
      .catch( error =>{
        console.log(error)
        this.error = error
      })
      .finally(()=> this.loading = false)
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
