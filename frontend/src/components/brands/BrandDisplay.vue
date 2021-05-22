<template>
  <div>
    <v-col>
    <h2>{{brandInput.Name}}</h2>
    </v-col>
    <v-col v-if="brandInput.Image != 0">
      Illustration : 
      <v-img :eager="true" :src="imgSrc" max-width="500">
      </v-img>
    </v-col>
    <v-col>
      <h4>Description</h4>
      <pre>
      {{brandInput.Description}}
      </pre>
    </v-col>
    <v-row>
      <v-col col="2">
        Creation Year : {{ brandInput.CreationYear}}
      </v-col>
      <v-col col="2" v-if="brandInput.EndYear">
        Closing Year : {{ brandInput.EndYear}}
      </v-col>
      <v-col>
        <a :href="brandInput.WikiHref">
          Wiki 
        </a><br/>
        <a :href="brandInput.Href">
          Site
        </a>
        
      </v-col>
    </v-row>
  </div>
</template>

<script>

import ImagesService from '../../services/ImagesService'

export default {
  name: 'BrandDisplay',
  props: {'brandInput': Object },
  data: function(){
    return {
      imgSrc : "",
      imgID: 0
    }
  },
  mounted: function(){
    this.imgID = this.brandInput.Image
    console.log(ImagesService.getImagePath(this.imgID))
    ImagesService.getImagePath(this.imgID).then((result) =>{
      console.log(result)
      this.imgSrc = result
    })
    console.log("prout" + this.imgSrc)
  },
  updated: function(){
    
    if (this.brandInput.Image != this.imgID){
      console.log("New Image ! "+this.brandInput.Image)
      this.imgID = this.brandInput.Image
      ImagesService.getImagePath(this.brandInput.Image).then(result =>{
        console.log(result)
        this.imgSrc = result
        console.log("image Source is : "+this.imgSrc);
      })


    }else{
      this.imgSrc = undefined
    }
  }
}
</script>
