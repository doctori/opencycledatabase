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
        if(imageID != undefined){
            return http.get('/images/'+imageID)
            .then(response =>{
                if(response.data != undefined){
                    console.log(response.data);
                    return "http://localhost:8081/"+response.data.Path
                }else{
                    return undefined
                }
                    
            
            });
        }
    }
}

export default new ImagesService();
