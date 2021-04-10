import http from "../common/http-common";

class UploadImagesService {
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
}

export default new UploadImagesService();
