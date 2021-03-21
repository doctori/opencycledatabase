<template>
  <div class="standard-create" id="standard-create">
    <h2> Standard Name </h2>
  <div id="standardName">
    Name : <MDBInput v-model="std.name" label="Standard Name"/>
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
    <div v-if="value.Type == 'bool'" class="form-check">
      <input class="form-check-input" type="checkbox" v-on:change=setFieldValue($event,key,value.Type) true-value="true" false-value="false" />
    </div>
    <div v-else-if="value.Type == 'int'">
      <MDBInput v-bind:id="key" class="std-input" type="number" v-bind:label="value.Name" v-on:change=setFieldValue($event,key,value.Type) />
    </div>
    <div v-else-if="value.Type == 'nested' || value.Type == 'nestedArray'">
      <select>
        <option v-for="nestedStandard in nestedStandards[key]" v-bind:key="nestedStandard.ID">
          {{nestedStandard.Name}}
        </option>
      </select>
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
  <div id="save-results" v-if="saved">
  </div>
  <div id="save-error" v-if="saveError">
    {{saveError}}
  </div>
</div>

</template>

<script>
import axios from 'axios'
import { MDBInput } from 'mdb-vue-ui-kit';
export default {
  name: 'StandardCreate',
  props: {'standards':Array},
  components:{
     MDBInput,
  },
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
      'saved': false,
      'saveError': null,
      'brands': [],
      'countryList':[],
      'loading':false,
      'stdDefintion':null,
      'nestedStandards': Map,
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
    getStandards(type){
      this.loading = true;
      axios
      .get("/standards/"+type.toLowerCase())
      .then(response =>(
        this.nestedStandards[type]=response.data
      ))
      .finally(()=>{
        this.loading=false
      })
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
        this.refreshNestedStd()
      })
    },
    refreshNestedStd(){
      for( const [,value] of Object.entries(this.stdDefintion)) {
        if (value.Type == "nested" || value.Type == "nestedArray"){
          this.getStandards(value.NestedType)
        }
      }
    },
    setFieldValue(value,field,type){
      value = value.target.value
      if (type == "int"){
        value = Number(value)
      }
      this.std[field]=value
    },
    submitStandard(){
      axios.post('/standards/'+this.std.type.toLowerCase(),this.std)
      .then(result => (
        this.std = result.data,
        this.saved = true
      ))
      .catch(error =>{
        console.log(error)
        this.saveError = error
      })
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
