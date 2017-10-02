<template>
  <div class="Main">
    <h2>Manage Standards</h2>
    <div class="standards">
        <div class="loading" v-if="loading">
            Loading ...
        </div>
        <div v-if="error" class="error">
            {{ error }}
        </div>
        <router-view></router-view>
          <ul v-if="standards && standards.length">
            <li v-for="standard in standards" :key="standard.ID">
               <router-link :to="`/standards/${standard.ID}`" >Aller à {{standard.Name}}</router-link>
               <!-- <p><strong>{{standard.Name}}</strong></p>-->
                <p>{{standard.Description}}</p>
                
            </li>
        </ul>
    <ul>
    </ul>
    </div>
    
  </div>

</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      loading: false,
      standards: null,
      error: null
    }
  },
  created () {
    // récupérer les données lorsque la vue est créée et
    // que les données sont déjà observées
    this.fetchData()
  },
  watch: {
    // appeler encore la méthode si la route change
    '$route': 'fetchData'
  },
  methods: {
    fetchData () {
      this.loading = true
      axios.get(`http://localhost:8080/standards`)
        .then(response => {
        // JSON responses are automatically parsed.
          this.standards = response.data
          this.loading = false
        })
        .catch(e => {
          this.errors.push(e)
        })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
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
