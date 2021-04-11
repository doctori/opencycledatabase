import http from "../common/http-common";

class ImagesService {
    upload(file,onUploadProgress) {
        let formData = new FormData();
        formData.append("file",file);
        return http.post("/images", formData, {
            headers: {
                "Content-Type": "multipart/form-data"
            },
            onUploadProgress
        });
    }
    getImages() {
        return http.get("/images");
    }
    getImagePath(imageID){
        var imagePath = "";
        http.get('/images/'+imageID)
        .then(response =>{
          console.log("image Path is "+response.data.Path)
          imagePath = "http://localhost:8081/"+response.data.Path
        });
        return imagePath;
  
    }
}

export default new ImagesService();
