<template>
<v-container>
  <v-row>
    <v-col>
      <h2>
        {{ $t('components.title') }}
      </h2>
    </v-col>
  </v-row>
  <v-row>
    <v-col cols="2">
      <v-autocomplete v-model="selectedType"
        v-on:change="setSelectedType(selectedType)"
        :items="typesList" 
        item-text="Type" 
        item-value="Type" 
        :label="$t('components.type')"
        dense
      ><template slot="item"  slot-scope="data">
        {{ $t('components.'+camelToSnakeCase(data.item.Type)) }}
        </template>
        <template slot="selection" slot-scope="data">
          {{ $t('components.'+camelToSnakeCase(data.item.Type)) }}
        </template>
      </v-autocomplete>
    </v-col>
    <v-col cols="2">
      <v-autocomplete v-model="selectedStandard"
        :items="standards" 
        v-on:change="setSelectedStandard(selectedStandard)"
        item-text="Name" 
        item-value="ID" 
        :label="$t('components.standard')"
        dense
      ></v-autocomplete>
    </v-col>
    <v-col cols="2">
      <v-autocomplete v-model="selectedBrand"
        :items="brands" 
        v-on:change="setselectedBrand(selectedBrand)"
        item-text="Name" 
        item-value="ID" 
        :label="$t('messages.brand')"
        dense
      ></v-autocomplete>
    </v-col>
    <v-col cols="1">
    </v-col>
    <v-col cols="1" >
      <v-btn 
      id="edit"
      elevation="4"
      v-on:click="search()"
      >
        {{editMessage}}
      </v-btn>
    </v-col>
      <v-col cols="1" >
        <v-btn 
        id="create"
        elevation="4"
        v-on:click="changeCreateMode()"
        >
        {{ $t('messages.create') }}
        </v-btn>
   </v-col>
  </v-row>
  <component-display v-for="(component) in components" :key="component.ID"
    :componentInput="component" 
  />
  <component-edit 
    :componentInput="selectedComponent" 
    :typeInput="selectedType" 
    :standardInput="selectedStandard"
    :brandInput="selectedBrand"
    v-if="componentEdit" 
  />


</v-container>
</template>

<script>
import http from "../../common/http-common";
import ComponentDisplay from "./ComponentDisplay";
import ComponentCreate from './ComponentCreate';
import UtilService from '../../services/UtilService';
export default {
  name: 'ComponentsMain',
  components:{
    'component-display': ComponentDisplay,
    'component-edit': ComponentCreate,
    },
  data : function(){
    return {
      typesList: [],
      standards: [],
      brands: [],
      components : [],
      componentID : '',
      selectedType: Object,
      selectedComponent : Object,
      selectedStandard: Object,
      selectedBrand: Object,
      componentDisplay: true,
      componentEdit: false,
      editMessage: this.$t('messages.search')
    }
  },
  mounted(){
    // retrieve the typesList
    http.get("/standards")
    .then(response =>{
      this.typesList = response.data
    })
    // retrieve components (should we???)
    http.get("/components")
    .then(response => {
      this.components = response.data
    });
    // retrieve Brands
    http.get("/brands")
    .then(response=>{
      this.brands = response.data
    });
  },
  methods: {
    camelToSnakeCase(str){
      return UtilService.camelToSnakeCase(str);
    },
    search(){
      console.log("Filter is standard : ["+this.selectedStandard+"] ")
      console.log("Filter is Brand : ["+this.selectedBrand+"] ")
      http.get("/components",{
        params: {
            standard: this.selectedStandard,
            brand: this.selectedBrand,
          }})
      this.componentEdit = false
  
    },
    changeEditMode(){
      this.componentDisplay = ! this.componentDisplay
      this.componentEdit = ! this.componentEdit
      if (this.componentEdit){
        this.editMessage = this.$t('messages.view')
      }else{
        this.editMessage = this.$t('messages.edit')
      }
    },
    changeCreateMode(){
      // reset selected Brand
      this.selectedComponent = new(Object)
      this.componentEdit = true
      this.componentDisplay = !this.componentEdit
      
    },
    setSelectedBrand(selectedBrand){
      this.selectedBrand = selectedBrand
    },
    setSelectedType(selectedType){
      this.selectedType = selectedType
      this.selectedStandard = 0
      http.get("/standards/"+selectedType.toLowerCase())
      .then(response =>{
        this.standards = response.data
      })
    },
    setSelectedStandard(selectedStandard){
      this.selectedStandard = selectedStandard
      http.get("/components",{
        params:{
          standard: selectedStandard
        }
      })
      .then(response => {
        this.components = response.data
      })
    },
    setselectedComponent(componentID){
      http.get("/components/"+componentID)
      .then(response => {
        this.selectedComponent = response.data
      })
      // TODO : catch errors
     },
     setselectedBrand(selectedBrand){
       this.selectedBrand=selectedBrand
     }
  }
}

</script>
