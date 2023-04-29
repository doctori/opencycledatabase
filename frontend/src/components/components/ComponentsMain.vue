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
          v-on:change="setSelectedType()"
          :items="typesList" 
          item-title="Type" 
          item-value="Type" 
          :label="$t('components.type')"
          dense
        ><template slot="item"  slot-scope="data">
          {{ $t('types.'+camelToSnakeCase(data.item.Type)) }}
          </template>
          <template slot="selection" slot-scope="data">
            {{ $t('types.'+camelToSnakeCase(data.item.Type)) }}
          </template>
        </v-autocomplete>
      </v-col>
      <v-col cols="2">
        <v-autocomplete v-model="selectedStandard"
          :items="standards" 
          item-title="Name" 
          item-value="ID" 
          :label="$t('components.standard')"
          dense
        ></v-autocomplete>
      </v-col>
      <v-col cols="2">
        <v-autocomplete v-model="selectedBrand"
          :items="brands" 
          v-on:change="setSelectedBrand(selectedBrand)"
          item-title="Name" 
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
        :loading="loading"
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
    <component-display 
      v-on:edit-component="setEditMode"
      v-for="(component) in components" 
      :key="component.ID"
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
import BackendApiClient from '../../services/BackendApiClient';
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
      selectedType: '',
      selectedComponent : Object,
      selectedStandard: -1,
      selectedBrand: "",
      componentDisplay: true,
      componentEdit: false,
      loading: false,
      editMessage: this.$t('messages.search')
    }
  },
  mounted(){
    // retrieve the typesList
    BackendApiClient.get("/standards")
    .then(response =>{
      this.typesList = response.data;
      this.setSelectedType();
    })
    // retrieve components (should we???)
    BackendApiClient.get("/components")
    .then(response => {
      this.components = response.data
    });
    // retrieve Brands
    BackendApiClient.get("/brands")
    .then(response=>{
      this.brands = response.data
    });
  },
  watch: {
    selectedType(val) {
      console.log("Updated Type is "+ val)
      BackendApiClient.get("/standards/"+this.selectedType.toLowerCase())
        .then(response =>{
          this.standards = response.data
        })
    }
  },
  methods: {
    camelToSnakeCase(str){
      return UtilService.camelToSnakeCase(str);
    },
    search(){
      var filters = {}
      console.log("Selected Type is "+ this.selectedType)
      if (this.selectedType != -1 && this.selectedStandard == 0){
        filters["type"] = this.selectedType;
      }
      if (this.selectedStandard != 0){
        filters["standard"] = this.selectedStandard;
      }
      if (this.selectedBrand != ""){
        filters["brand"] = this.selectedBrand
      }
      console.log("Filter is Type : ["+filters["type"]+"]");
      console.log("Filter is standard : ["+this.selectedStandard+"] ");
      console.log("Filter is Brand : ["+this.selectedBrand+"] ");
      this.loading = true
      BackendApiClient.get("/components",{
        params: filters
      }).then(response => {
        if (response.length != 0){
          this.components = response.data;
          this.componentDisplay = true;
          this.componentEdit = false;
        }
        this.loading = false;
      })


      this.componentEdit = false
    },
    setEditMode(component){
      console.log(component);
      this.selectedComponent = component;
      this.changeEditMode();
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
    setSelectedType(){
        this.selectedStandard = 0
        console.log("Standard selected"+this.selectedType);
        BackendApiClient.get("/standards/"+this.selectedType.toLowerCase())
        .then(response =>{
          this.standards = response.data
        })
    },
    setSelectedStandard(selectedStandard){
      this.selectedStandard = selectedStandard
      BackendApiClient.get("/components",{
        params:{
          standard: selectedStandard
        }
      })
      .then(response => {
        this.components = response.data
      })
    },
  }
}

</script>
