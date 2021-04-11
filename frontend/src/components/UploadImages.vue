<template>
  <div>
    <div v-if="currentFile">
      <div>
        <v-progress-linear
          v-model="progress"
          color="light-blue"
          height="25"
          reactive
        >
          <strong>{{ progress }} %</strong>
        </v-progress-linear>
      </div>
    </div>
      <v-col>
        <v-file-input
          show-size
          label="Image input"
          @change="upload"
        ></v-file-input>
      </v-col>

    <v-alert v-if="message" border="left" color="blue-grey" dark>
      {{ message }}
    </v-alert>

  </div>
</template>

<script>
import UploadImage from "../services/UploadImagesServices";

export default {
    name: 'upload-images',
    data(){
        return {
            currentFile: undefined,
            progress: 0,
            message: "",
            fileInfos: []
        };
    },
    mounted(){
        UploadImage.getImages().then(response=>{
            this.fileInfos = response.data;
        })
    },
    methods: {
        upload(file) {
            this.currentFile = file;
            if(!this.currentFile){
                this.message = "Please select a file";
                return;
            }
            this.message = "";
            UploadImage.upload(this.currentFile, (event) => {
                this.progress = Math.round((100*event.loaded)/event.total);
            })
            .then((response) => {
                this.message = response.data.message;
                this.$emit('image-uploaded',response.data.ID)
            })
            .catch(()=> {
                this.progress = 0;
                this.message = "Could not upload the file";
                this.currentFile = undefined;
            })
        }
    }
};
</script>
