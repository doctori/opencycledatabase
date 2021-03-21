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
  <!-- let's display the common fields for all standards -->
  <!-- let's list all the countries !! -->
  <div id="standardCountry">
    Country :
    <select v-model="std.country" >
      <option label="none"></option>
      <option v-for="country in countryList" v-bind:key="country.alpha3Code">
        {{country.name}}
      </option>
    </select>
  </div>
  <!-- let's list the known brands -->
  <div id="standardBrand">
    Brand :
    <select v-model="std.brand" >
      <option label="none"></option>
      <option v-for="brand in brands" v-bind:key="brand.Name">
        {{brand.Name}}
      </option>
    </select>
  </div>
  <div v-if="loading">
    Loading ...
  </div>
  <div v-for="(value,key) in stdDefintion" v-bind:key="key">
    <label v-bind:id="key" class="std-field">
      {{key}}
    </label>
    <div v-if="value.Type == 'bool'">
      <input type="checkbox" v-on:change=setFieldValue($event,key) true-value="true" false-value="false" >
    </div>
    <div v-else>
      <input v-bind:id="key" class="std-input" v-bind:key="key" v-on:change=setFieldValue($event,key)>
    </div>
  </div>
  <div id="result">
    {{ std }}
  </div>
  <div v-if="error">
    HAAAAAAAAAA {{error}}
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
        'type':'',
        // TODO : get Country list
        'country':'',
        // TODO : get brand
        'brand':''
      },
      'brands': [],
      'countries':[],
      'loading':false,
      'stdDefintion':null,
      'error': null,
    }
  },
  mounted: function (){
    axios
        .get("https://restcountries.eu/rest/v2/all")
        .then(response => (
          this.countryList = response.data
        ))
    axios
      .get("/brands")
      .then(response => (
        this.brands = response.data
      ))
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
    setFieldValue(value,field){
      this.std[field]=value.target.value
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
