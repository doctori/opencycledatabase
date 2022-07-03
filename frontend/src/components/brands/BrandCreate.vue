<template>
<v-card>
  <v-form fluid>
    <div class="brand-create" id="standard-create">
      <h2> Brand : {{brand.name}}</h2>
    <v-col id="brandName">
      <v-text-field v-model="brand.Name" label="Brand Name" required>Name</v-text-field>
    </v-col>
    
    <!-- let's display the common fields for all standards -->
    <!-- let's list all the countries !! -->
    <v-row>
      <v-col id="brandCountry">
        Country :
        <v-autocomplete v-model="brand.Country" label="country name" :items="countryList" item-text="name" item-value="alpha3Code" >
        </v-autocomplete>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Description
          <v-textarea  v-model="brand.Description" label="description">
          </v-textarea>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Creation Year
          <v-text-field type="number" v-model.number="brand.CreationYear" label="Creation Year">
          </v-text-field>
      </v-col>
      <v-col>
        End Year
          <v-text-field type="number" v-model.number="brand.EndYear" label="End Year">
          </v-text-field>
      </v-col>
      <v-col>
        Image
        <v-img :eager="true" :src="imgSrc" max-width="200"></v-img>
        <upload-image v-on:image-uploaded="setBrandImage"></upload-image>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        Wikipedia link
          <v-text-field  v-model="brand.WikiHref" label="Wiki Link">
          </v-text-field>
      </v-col>
      <v-col>
        Official website link
          <v-text-field v-model="brand.Href" label="End Year">
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
import http from '../../common/http-common'
import BackendApiClient from '../../services/BackendApiClient'
import UploadImage from './../UploadImages'
import ImagesService from '../../services/ImagesService'
export default {
  name: 'BrandCreate',
  components:{
    'upload-image' : UploadImage,
  },
  props: {
    brand: Object
  },
  data : function(){
    return {
      'saved': false,
      'saveError': null,
      'countryList':[],
      'loading':false,
      'error': null,
      'imgSrc': "",
      'imgID': 0,
    }
  },  
  updated: function(){

   
      ImagesService.getImagePath(this.brand.Image).then((result) =>{
      console.log(result)
      this.imgSrc = result
    })
  },
  mounted: function (){
    http
        .get("https://restcountries.eu/rest/v2/all")
        .then(response => (
          this.countryList = response.data
        ))
      if (this.brand.Image != 0){
         ImagesService.getImagePath(this.imgID).then((result) =>{
          console.log(result)
          this.imgSrc = result
        })
      }
  },
  methods: {
    includeFields(field){
      return !this.ignoredFields.includes(field)

    },
    setBrandImage(image){
      console.log("this brand Image ID is "+image)
      this.brand.Image=image
    },
    submitBrand(){
      if (this.brand.ID != 0 && this.brand.ID != undefined){
        console.log("we'll update the Brand "+this.brand.Name)
        BackendApiClient.post('/brands/'+this.brand.ID,this.brand)
        .then(result => (
          this.brand = result.data,
          this.saved = true
        ))
        .catch(error =>{
          console.log(error)
          this.saveError = error
        })
      }else{
        BackendApiClient.post('/brands',this.brand)
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
