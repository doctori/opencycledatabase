<template>
<v-card>
  <v-form fluid>
    <v-row align="center"  >
      <v-col cols="12" id="standardType">
      Standard Type : 
      <v-autocomplete v-model="std.type" v-on:change="getStandardDefintion(std.type)"
        :items="stdTypes" outlined dense chips
      ></v-autocomplete>
    </v-col>
    </v-row>
    <div class="standard-create" id="standard-create">
      <h2> Standard Name </h2>
    <v-col id="standardName">
      <v-text-field v-model="std.name" label="Standard Name" required>Name</v-text-field>
    </v-col>
    
    <!-- let's display the common fields for all standards -->
    <!-- let's list all the countries !! -->
    <v-row>
      <v-col id="standardCountry">
        Country :
        <v-autocomplete v-model="std.country" label="country name" :items="countryList" item-text="name" item-value="alpha3Code" >
        </v-autocomplete>
      </v-col>

      <!-- let's list the known brands -->
      <v-col id="standardBrand">
        Brand :
        <v-autocomplete v-model="std.brand" label="brand name" :items="brands" item-text="Name" item-bind="ID">
        </v-autocomplete>
      </v-col>
    </v-row>
    <div v-if="loading">
      Loading ...
    </div>
    <div v-else v-for="(value,key) in stdDefintion" v-bind:key="key">
      <v-row v-if="value.Type == 'bool'" >
        <v-col>
        <v-switch v-model="std[value.Name]" class="form-check" :label=key>
        </v-switch>
        </v-col>
      </v-row>
      <v-row v-else-if="value.Type == 'int'">
        <v-col  sm="6" cols="10">
          <v-text-field  v-model="std[value.Name]" :suffix="value.Unit" :label=key>
          </v-text-field>
        </v-col>
      </v-row>
      <v-row v-else-if="value.Type == 'nested' || value.Type == 'nestedArray'">
        <v-select v-model="std[value.Name]"  :items="nestedStandards[key]" item-text="Name" item-bind="ID" :label=key>
        </v-select>
      </v-row>
      <v-row v-else>
        <v-col sm="6" cols="10">
          <v-text-field v-model="std[value.Name]" :label=key >
          </v-text-field>
        </v-col>
      </v-row>
    </div>
    <div id="result">
      {{ this.std }}
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
  </v-form>
</v-card>
</template>

<script>
import axios from 'axios'
export default {
  name: 'StandardCreate',
  props: {'standards':Array},
  components:{
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
      'stdTypes': [],
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
      // build standard type list : 
      this.standards.forEach(std => {
        this.stdTypes.push(std.Type)
      });

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

      //value = value.target.value
      console.log(value)
      console.log(field)
      console.log(type)
      if (type == "int"){
        value = Number(value)
      }
      this.std[field]=value
      console.log(this.std)
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
