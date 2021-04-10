<template>
<v-card>
  <v-form fluid>
    <div class="brand-create" id="standard-create">
      <h2> Brand : {{brand.name}}</h2>
    <v-col id="brandName">
      <v-text-field v-model="brand.name" label="Brand Name" required>Name</v-text-field>
    </v-col>
    
    <!-- let's display the common fields for all standards -->
    <!-- let's list all the countries !! -->
    <v-row>
      <v-col id="brandCountry">
        Country :
        <v-autocomplete v-model="brand.country" label="country name" :items="countryList" item-text="name" item-value="alpha3Code" >
        </v-autocomplete>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Description
          <v-textarea  v-model="brand.description" label="description">
          </v-textarea>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Creation Year
          <v-text-field type="number" v-model.number="brand.creationYear" label="Creation Year">
          </v-text-field>
      </v-col>
      <v-col>
        End Year
          <v-text-field type="number" v-model.number="brand.endYear" label="End Year">
          </v-text-field>
      </v-col>
      <v-col>
        Image
        <upload-image></upload-image>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Wikipedia link
          <v-text-field  v-model="brand.wikiHref" label="Wiki Link">
          </v-text-field>
      </v-col>
      <v-col>
        Official website link
          <v-text-field v-model="brand.href" label="End Year">
          </v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-btn block elevation="2" x-large key="submit" v-on:click="submitBrand()">
          submit
        </v-btn>
      </v-col>
    </v-row>
    <div id="result">
      {{ this.brand }}
    </div>

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
import UploadImage from './UploadImages'
export default {
  name: 'BrandCreate',
  components:{
    'upload-image' : UploadImage,
  },
  data : function(){
    return {
      'brand':{
        'name':'',
        'type':'',
        // TODO : get Country list
        'country':'',
        // TODO : get brand
        'brand':''
      },
      'saved': false,
      'saveError': null,
      'countryList':[],
      'loading':false,
      'error': null,
    }
  },
  mounted: function (){
    axios
        .get("https://restcountries.eu/rest/v2/all")
        .then(response => (
          this.countryList = response.data
        ))
  },
  methods: {
    includeFields(field){
      return !this.ignoredFields.includes(field)

    },
    submitBrand(){
      axios.post('/brands',this.brand)
      .then(result => (
        this.brand = result.data,
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
.brand-field {
  margin: 0 20px 0;
}
h3 {
  margin: 40px 0 0;
}
a {
  color: #42b983;
}
</style>
