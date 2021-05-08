<template>
<v-container>
  <v-row>
    <v-col cols="1">
      Types : 
    </v-col>
    <v-col cols="6">

      <v-autocomplete v-model="selectedType"
        v-on:change="setSelectedType(selectedType)"
        :items="standardTypes" 
        item-text="Type" 
        item-value="Type" 
        :label="$t('components.type')"
        dense
      ><template slot="item"  slot-scope="data">
        {{ $t('types.'+camelToSnakeCase(data.item)) }}
        </template>
        <template slot="selection" slot-scope="data">
          {{ $t('types.'+camelToSnakeCase(data.item)) }}
        </template>
      </v-autocomplete>
    </v-col>
    <v-col cols="1" >
      <v-btn 
      id="edit"
      elevation="4"
      v-on:click="changeEditMode()"
      :disabled="selectedType===''"
      >
        {{ $t('messages.' + editMessage ) }}
      </v-btn>
    </v-col>
  </v-row>
  <div id="standard-display" v-if="displayMode && standards.length != 0" >
    <standard-display v-on:edit-standard="setEditMode" v-for="(standard) in standards" :standardInput="standard" :key="standard.ID"  />
  </div>
  <div id="empty-display" v-else-if="displayMode">
    {{ $t('messages.empty_set')}}
  </div>
  <standard-edit :standardTypeInput="selectedType" :standardInput="selectedStandard" v-if="editMode"/>

</v-container>
</template>

<script>
import http from "../../common/http-common";
import StandardDisplay from './StandardDisplay';
import StandardCreate from './StandardCreate';
import UtilService from '../../services/UtilService';
export default {
  name: 'StandardsMain',
  components:{
    'standard-display': StandardDisplay,
    'standard-edit': StandardCreate,
    },
  data : function(){
    return {
      standards : [],
      standardTypes: [],
      standardID : '',
      selectedType: "",
      selectedStandard : Object,
      displayMode: true,
      editMode: false,
      editMessage: "create"
      
    }
  },
  mounted(){
    http.get("/standards")
    .then(response => {
      // build standard type list : 
      response.data.forEach(std => {
        this.standardTypes.push(std.Type)
      });
    }
 

    )
  },
  methods: {
    camelToSnakeCase(str){
      return UtilService.camelToSnakeCase(str);
    },
    changeEditMode(){
      this.editMode = ! this.editMode
      this.displayMode = ! this.editMode
      if (this.editMode){
        this.editMessage = "view"
      }else{
        this.editMessage = "edit"
      }
    },
    setEditMode(standard){
      console.log(standard);
      this.selectedStandard = standard;
      this.changeEditMode();
    },
    changeCreateMode(){
      // reset selected Brand
      this.editMode = true
      this.displayMode = !this.brandEdit
      
    },
    setSelectedType(selectedType){
      this.selectedStandard = 0
      http.get("/standards/"+selectedType.toLowerCase())
      .then(response =>{
        this.standards = response.data
        console.log(this.standards.length);
      })

    },
    setSelectedStandard(standardID){
      http.get("/standards/"+standardID)
      .then(response => {
        this.selectedStandard = response.data

      })
      // TODO : catch errors
     }
  }
}
  
</script>
