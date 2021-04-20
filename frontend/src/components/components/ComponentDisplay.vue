<template>
  <div>
    <v-col v-if="componentInput.Image != 0">
      Illustration : 
      <v-img v-if="imgSrc" :src="imgSrc" max-width="350" :eager="true">
      </v-img>
    </v-col>
    <v-col>
      <h4>Description</h4>
      <pre>
      {{componentInput.Description}}
      </pre>
    </v-col>
    <v-row>
      <v-col col="2">
        Creation Year : {{ componentInput.CreationYear}}
      </v-col>
      <v-col col="2" v-if="componentInput.EndYear">
        Closing Year : {{ componentInput.EndYear}}
      </v-col>
      <v-col>
        <a :href="componentInput.WikiHref">
          Wiki 
        </a><br/>
        <a :href="componentInput.Href">
          Site
        </a>
        
      </v-col>
    </v-row>
  </div>
</template>

<script>

import ImagesService from '../../services/ImagesService'
export default {
  name: 'ComponentDisplay',
  props: {'componentInput': Object },
  data: function(){
    return {
      imgSrc : "",
      imgID: 0
    }
  },
  mounted: function(){
    this.imgID = this.componentInput.Image
    this.imgSrc = ImagesService.getImagePath(this.imgID)
  },
  updated: function(){
    console.log("updated")
    if (this.componentInput.Image != this.imgID){
      this.imgID = this.componentInput.Image
      this.imgSrc = ImagesService.getImagePath(this.imgID)
    }else{
      this.imgSrc = undefined
    }
  }
}
</script>
