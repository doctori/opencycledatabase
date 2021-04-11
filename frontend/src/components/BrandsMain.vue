<template>
<div>
  <v-col id="Brands">
    Brands : 
    <v-autocomplete v-model="brandID"
      :items="brands" 
      v-on:change="setSelectedBrand(brandID)"
      item-text="Name" 
      item-value="ID" 
      outlined 
      dense
    ></v-autocomplete>
  </v-col>
  <v-col>
    <brand-display :brand="selectedBrand"/>
  </v-col>
</div>
</template>

<script>
import http from "../common/http-common";
import BrandDisplay from "./BrandDisplay";

export default {
  name: 'BrandsMain',
  components:{
    'brand-display': BrandDisplay,
  },
  data : function(){
    return {
      brands : [],
      selectedBrand : Object,
    }
  },
  mounted(){
    http.get("/brands")
    .then(response => {
      this.brands = response.data
    }

    )
  },
  methods: {
    setSelectedBrand(brandID){
      http.get("/brands/"+brandID)
      .then(response => {
        this.selectedBrand = response.data
      })
      // TODO : catch errors
     }
  }
}
  
</script>
