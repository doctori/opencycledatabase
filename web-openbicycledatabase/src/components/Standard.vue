<template>
    <div class="standard">
        PROUT ? 
        <div class="loading" v-if="loading">
            Loading ...
        </div>
        <div v-if="error" class="error">
            {{ error }}
        </div>
          <ul v-if="standard">
            <h3>{{standard.Name}}</h3>
            <li>              
                <p><strong>{{standard.Country}}</strong></p>
            </li>
            <li>
                <p>{{standard.Code}}</p>
            </li>
            <li>
                <p>{{standard.Description}}</p>
            </li>
        </ul>
    </div>

</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      loading: false,
      standard: null,
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
      axios.get('http://localhost:8080/standards/' + this.$route.params.id)
        .then(response => {
        // JSON responses are automatically parsed.
          this.standard = response.data
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
