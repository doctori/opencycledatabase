<template>
  <div>
    <v-row>
      <v-col cols="3">
        <h2>{{ componentInput.Name }}</h2>
      </v-col>
      <v-col cols="2">
        <h4>{{ componentInput.Brand.Name }}</h4>
      </v-col>
      <v-col>
        <v-img v-if="imgSrc" :src="imgSrc" max-width="350" :eager="true">
        </v-img>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <h4>Description</h4>
        <pre>
        {{componentInput.Description}}
        </pre>
      </v-col>
      <v-col cols="2">
        <v-btn
          id="edit"
          v-on:click="$emit('edit-component',componentInput)">
          {{ $t('messages.edit')}}
        </v-btn>
        <v-btn
          id="delete"
          v-on:click="deleteComponent()">
          {{ $t('messages.delete')}}
        </v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import http from "../../common/http-common";
import ImagesService from '../../services/ImagesService'
export default {
  name: 'ComponentDisplay',
  props: {
    'componentInput': Object 
  },
  data: function(){
    return {
      imgSrc : "",
      imgID: 0
    }
  },
  mounted: function(){
    console.log(this.componentInput)
    this.imgID = this.componentInput.Image
    if (this.imgID != undefined){
      ImagesService.getImagePath(this.imgID).then(result =>{
        this.imgSrc = result
      })
    }
    
    
  },
  updated: function(){
    if (this.componentInput.Image != this.imgID){
      this.imgID = this.componentInput.Image
      ImagesService.getImagePath(this.imgID).then(result =>{
        if(result != undefined){
          this.imgSrc = result
        }
      })
    }else{
      this.imgSrc = undefined
    }
  },
  methods: {
    deleteComponent(){
      http.delete('/components/'+this.componentInput.ID)
      this.$emit('delete-component',this.componentInput)
      this.componentInput = null
    }
  }

}
</script>
