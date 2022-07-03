import http from "../common/http-common";
class BackendApiClient {
    rootURL = "/api"
    get(url,config){
        return http.get(this.rootURL  + url,config)
    }
    post(url,data,config){
        return http.post(this.rootURL + url,data,config)
    }
    delete(url,config) {
        return http.delete(this.rootURL + url,config)
    }
}
export default new BackendApiClient();
