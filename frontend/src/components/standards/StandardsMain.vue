<template>
<v-container>
  <v-row>
    <v-col cols="1">
      Standards : 
    </v-col>
    <v-col cols="6">

      <v-autocomplete v-model="standardID"
        :items="standards" 
        v-on:change="setSelectedStandard(standardID)"
        item-text="Name" 
        item-value="ID" 
        dense
      ></v-autocomplete>
    </v-col>
    <v-col cols="1" >
      <v-btn 
      id="edit"
      elevation="4"
      v-on:click="changeEditMode()"
      >
        {{editMessage}}
      </v-btn>
    </v-col>
    <v-col cols="1">
      <v-btn v-on:click="changeCreateMode()">
        Create
      </v-btn>
   </v-col>
  </v-row>
  <standard-display :standardInput="selectedStandard" v-if="displayMode && Object.keys(selectedStandard).length != 0 "/>
  <standard-edit :standardTypes="standardTypes" :standard="selectedStandard" v-if="editMode" />

</v-container>
</template>

<script>
import http from "../../common/http-common";
import StandardDisplay from './StandardDisplay';
import StandardCreate from './StandardCreate';
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
      selectedStandard : Object,
      displayMode: true,
      editMode: false,
      editMessage: "edit"
      
    }
  },
  mounted(){
    http.get("/standards")
    .then(response => {
      this.standards = response.data
      console.log(this.standards)
      // build standard type list : 
      this.standards.forEach(std => {
        this.standardTypes.push(std.Type)
      });
    }
 

    )
  },
  methods: {
    changeEditMode(){
      this.displayMode = ! this.displayMode
      this.editMode = ! this.editMode
      if (this.editMode){
        this.editMessage = "view"
      }else{
        this.editMessage = "edit"
      }
    },
    changeCreateMode(){
      // reset selected Brand
      this.selectedStandard = new(Object)
      this.editMode = true
      this.displayMode = !this.brandEdit
      
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
