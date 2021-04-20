<template>
<v-card>
  <v-form fluid>
    <v-row >
      <v-col cols="2">
        <h2> {{ $t('components.name') }} </h2>
      </v-col>
      <v-col cols="9">
        <v-text-field v-model="componentInput.name"
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
        <v-textarea v-model="componentInput.description"
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
      {{ this.componentInput }}
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
//import http from '../../common/http-common'
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
    standardList: Array
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
  beforeUpdate: function(){
    // we update the selected type if needed
    console.log(this.typeInput)
    this.componentInput.Type = this.typeInput
  },
  mounted: function (){
  },
  methods: {
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
