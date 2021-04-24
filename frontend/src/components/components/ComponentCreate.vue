<template>
<v-card>
  <v-form fluid>
    <v-row >
      <v-col cols="2">
        <h2> {{ $t('components.name') }} </h2>
      </v-col>
      <v-col cols="9">
        <v-text-field v-model="cpn.name"
                      :label="$t('components.name')" required>
          {{ $t('components.name') }}
        </v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="2">
        {{ $t('components.description') }}
      </v-col >
      <v-col>
        <v-textarea v-model="cpn.description"
                      :label="$t('components.description')">
            {{ $t('components.description') }}
        </v-textarea>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-btn block elevation="2" x-large key="submit" v-on:click="submitComponent()">
          {{ $t('messages.submit') }}
        </v-btn>
      </v-col>
    </v-row>
    
    <div id="result">
      {{ this.cpn }}
    </div>

    <div id="save-results" v-if="saved">
    </div>
    <div id="save-error" v-if="saveError">
      {{saveError}}
    </div>
  </v-form>
</v-card>
</template>

<script>
import http from '../../common/http-common'
//import UploadImage from './../UploadImages'
//import ImageService from '../../services/ImagesService'
import UtilService from '../../services/UtilService'
export default {
  name: 'ComponentCreate',
  components:{
//    'upload-image' : UploadImage,
  },
  props: {
    componentInput: Object,
    typeInput: String,
    standardInput: Number,
    standardList: Array,
    brandInput: Number
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
      'cpn': Object
    }
  },  
  updated: function(){
    // we update the selected type if needed
    console.log(this.typeInput)
    this.cpn.Type = this.typeInput
    this.cpn.Standard = this.standardInput
    this.updateBrand(this.brandInput)
  },
  mounted: function (){
    this.cpn = this.componentInput
  },
  methods: {
    updateBrand(newBrandID){
      if (newBrandID != this.cpn.Brand.ID){
        http.get('/brands/'+newBrandID)
        .then(result => {
          this.cpn.Brand = result.data
        });
      }
    },
    updateStandard(newStandardID){
      if (newStandardID != this.cpn.Standard.ID){
        http.get('/standards/+newStandardID')
        .then(result => {
          this.cpn.Standard = result.data
        });
      }
    },
    submitComponent(){
    http
      .post('/components/',this.cpn)
      .then(result => (
        this.cpn = result.data,
        this.saved = true
      ))
      .catch(error =>{
        console.log(error)
        this.saveError = error
      })
    },
    camelToSnakeCase(str){
      return UtilService.camelToSnakeCase(str);
    },
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
