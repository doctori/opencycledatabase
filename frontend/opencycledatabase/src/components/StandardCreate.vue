<template>
  <div class="standard-create" id="standard-create">
    <h2> Standard Name </h2>
  <div id="standardName">
    Name : <input v-model="std.name" placeholder="standard Name">
  </div>
  <div id="standardType">
    Standard Type : 
    <select v-model="std.type" v-on:change="getStandardDefintion(std.type)">
      <option v-for="standard in standards" v-bind:key="standard.Type">
        {{standard.Type}}
      </option>
    </select>
  </div>
  <div>
    <div v-if="loading">
      Loading ...
    </div>
    <div v-for="(value,key) in stdDefintion" v-bind:key="key">
      <div v-bind:id="key" v-if="includeFields(key)">
        <label v-bind:id="key" class="std-field">
          {{key}}
        </label>
        <input v-bind:id="key" class="std-input" v-bind:key="key" v-model="form.parent_id[n]">
      </div>
    </div>
  </div>
    <div id="result">
      {{ std.name }} : {{ std.type }}
    </div>
    <div v-if="error">
      HAAAAAAAAAA {{error}}
    </div>
    <div v-if="stdDefintion" id="type def">
      Definition : 
      {{stdDefintion}}
    </div>
    <button key="submit" v-on:click="submitStandard()">
      submit
    </button>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'StandardCreate',
  props: {'standards':Array},
  data : function(){
    return {
      'std':{
        'name':'',
        'type':''
      },
      'loading':false,
      'stdDefintion':null,
      'error': null,
      'ignoredFields':[
        'ID',
        'Name',
        'CreatedAt',
        'UpdatedAt',
        'DeletedAt'
      ]
    }
  },
  methods: {
    includeFields(field){
      return !this.ignoredFields.includes(field)

    },
    getStandardDefintion(type){
      this.loading = true;
      axios
      .get("/standards/"+type.toLowerCase(),{
        params: {
          struct_only: true
        }
      })
      .then(response => (
        this.stdDefintion = response.data,
        this.error = false
         ))
      .catch( error =>{
        console.log(error)
        this.error = error.message
        this.stdDefintion = null
      })
      .finally(()=>{
        this.loading = false
      })
    },
    submitStandard(){
      console.log(this.std)
    }
  }

}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.std-field {
  margin: 0 20px 0;
}
h3 {
  margin: 40px 0 0;
}
a {
  color: #42b983;
}
</style>
