<template>
<div>
    <v-main>
      <v-container >
        <v-row>
          <div v-if="loading" class="loading">
            Loading...
          </div>
          <div v-if="error" class="error">
            {{ error.message }}
          </div>
          <v-col v-if="standards && !loading" class="content">
            <standard-display :standards="standards"/>
          </v-col>
      </v-row>
      </v-container>
      
      <v-container>
        <v-row>
          <v-col >
            <v-btn id="createStandard" v-on:click="enableCreateStandard()">
              Create Standard
            </v-btn>
          </v-col>
        </v-row>
        <v-row>
          <v-col v-if="createStandard" class="create-standard">
            <standard-create :standards="standards"/>
          </v-col>
        </v-row>
      </v-container>
   </v-main>

</div>
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
      createStandard: false,
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
