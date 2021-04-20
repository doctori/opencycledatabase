<template>
<v-card>
  <v-form fluid>
    <v-row >
      <v-col>
        <h2> Standard :{{std.name}}</h2>
      </v-col>
    </v-row>
    <v-row class="standard-create" id="standard-create">
      <v-col id="standardName">
        <v-text-field v-model="std.name" label="Standard Name" required>Name</v-text-field>
      </v-col>
      <v-col id="standardType">
        Standard Type : 
        <v-autocomplete v-model="std.type" v-on:change="getStandardDefintion(std.type)"
          :items="standardTypes" outlined dense
        ></v-autocomplete>
      </v-col>
    </v-row>

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
          <v-text-field type="number" v-model.number="std[value.Name]" :suffix="value.Unit" :label=key>
          </v-text-field>
        </v-col>
      </v-row>
      <v-row v-else-if="value.Type == 'nested' || value.Type == 'nestedArray'">
        <v-col>
          <v-select 
            v-model="std[value.Name]"
            :items="nestedStandards[key]"
            :label=key
            >
          </v-select>
        </v-col>
      </v-row>
      <v-row v-else>
        <v-col sm="6" cols="10">
          <v-text-field v-model="std[value.Name]" :label=key >
          </v-text-field>
        </v-col>
      </v-row>
    </div>
    <v-row>
      <v-col>
        <v-btn block elevation="3" key="submit" v-on:click="submitStandard()">
          submit
        </v-btn>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="2">
      <h4>Result</h4>
      </v-col>
      <v-col cols="10">
        <pre>
            {{std}}
        </pre>
      </v-col>
    </v-row>
    <v-row>
      <div v-if="error">
        HAAAAAAAAAA {{error}}
      </div>
      <div id="save-results" v-if="saved">
      </div>
      <div id="save-error" v-if="saveError">
        {{saveError}}
      </div>
    </v-row>
  </v-form>
</v-card>
</template>

<script>
import http from '../../common/http-common';
export default {
  name: 'StandardCreate',
  props: {'standardTypes':Array},
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
    http
        .get("https://restcountries.eu/rest/v2/all")
        .then(response => (
          this.countryList = response.data
        ))
    http
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
      http
      .get("/standards/"+type.toLowerCase())
      .then(response =>{
        var thisNestedStandard = []
        response.data.forEach(element => {
          var standard = {}
          standard.value = element
          standard.text = element.Name
          thisNestedStandard.push(standard)
        });
        this.nestedStandards[type]=thisNestedStandard;
        console.log(thisNestedStandard);
        console.log("Adding to the type "+type);
      })
      .finally(()=>{
        this.loading=false
      })
    },
    getStandardDefintion(type){
      this.loading = true;
      http
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
    setNestedStandard(value,field,type){
      console.log(value)
      console.log(field)
      console.log(type)
      this.std[field]=value
    },
    submitStandard(){
      http
      .post('/standards/'+this.std.type.toLowerCase(),this.std)
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
