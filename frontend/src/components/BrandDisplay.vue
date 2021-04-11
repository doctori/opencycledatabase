<template>
  <div>
    <v-col v-if="brand.Image != 0">
      Illustration : 
      <v-img v-if="brand.ImagePath" :src="brand.ImagePath">
      </v-img>
    </v-col>
    <v-col>
      <h4>Description</h4>
      <pre>
      {{brand.Description}}
      </pre>
    </v-col>
  </div>
</template>

<script>
import http from '../common/http-common'
export default {
  name: 'BrandDisplay',
  props: {'brand': Object },
  data: function(){
    return {
      imgSrc : ""
    }
  },
  beforeUpdate: function(){
    console.log("updated")
    if (this.brand.Image != 0){
      http.get('/images/'+this.brand.Image)
      .then(response =>{
        console.log("image Path is "+response.data.Path)
        this.brand.ImagePath = response.data.Path  
      })
    }else{
      this.brand.ImagePath = undefined
    }
  }
}
</script>
