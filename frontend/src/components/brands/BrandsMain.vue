<template>
<v-container>
  <v-row>
    <v-col cols="1">
      Brands : 
    </v-col>
    <v-col cols="6">

      <v-autocomplete v-model="brandID"
        :items="brands" 
        v-on:change="setSelectedBrand(brandID)"
        item-text="Name" 
        item-value="ID" 
        dense
      ></v-autocomplete>
    </v-col>
    <v-col cols="2" >
      <v-btn 
      id="edit"
      elevation="4"
      v-on:click="changeEditMode()"
      >
        {{editMessage}}
      </v-btn>
    </v-col>
      <v-col cols="2" >
        <v-btn 
        id="create"
        elevation="4"
        v-on:click="changeCreateMode()"
        >
          {{ $t('messages.create') }}
        </v-btn>
   </v-col>
  </v-row>
  <brand-display :brandInput="selectedBrand" v-if="brandDisplay && Object.keys(selectedBrand).length != 0 "/>
  <brand-edit :brand="selectedBrand" v-if="brandEdit" />


</v-container>
</template>

<script>

import BackendApiClient from "../../services/BackendApiClient";
import BrandDisplay from "./BrandDisplay";
import BrandCreate from './BrandCreate';
export default {
  name: 'BrandsMain',
  components:{
    'brand-display': BrandDisplay,
    'brand-edit': BrandCreate,
    },
  data : function(){
    return {
      brands : [],
      brandID : '',
      selectedBrand : Object,
      brandDisplay: true,
      brandEdit: false,
      brandCreate: false,
      editMessage: "edit"
    }
  },
  mounted(){
    BackendApiClient.get("/brands")
    .then(response => {
      this.brands = response.data
    }

    )
  },
  methods: {
    changeEditMode(){
      this.brandDisplay = ! this.brandDisplay
      this.brandEdit = ! this.brandEdit
      if (this.brandEdit){
        this.editMessage = "view"
      }else{
        this.editMessage = "edit"
      }
    },
    changeCreateMode(){
      // reset selected Brand
      this.selectedBrand = new(Object)
      this.brandEdit = true
      this.brandDisplay = !this.brandEdit
      
    },
    setSelectedBrand(brandID){
      BackendApiClient.get("/brands/"+brandID)
      .then(response => {
        this.selectedBrand = response.data
      })
      // TODO : catch errors
     }
  }
}
  
</script>
