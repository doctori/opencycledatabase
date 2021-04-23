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
        return http.get('/images/'+imageID)
        .then(response =>{
          return "http://localhost:8081/"+response.data.Path
          
        });
        
        
    }
}

export default new ImagesService();
