package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO : move to gorilla mux ...

// Resource will define the generic API resource methods
// NonJSONResource hold the files images and other objects
type NonJSONResource interface {
	Get(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{})
	Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{})
	Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{})
}

// NonJSONData represent any content that is not JSon, mostly images
type NonJSONData interface {
	GetContentType() string
	GetContentLength() string
	GetContent() []byte
}

// API is the generic "API" model
type API struct{}

func (api *API) splitPath(path string, resourceType string) (id primitive.ObjectID, adj string) {
	id = primitive.NilObjectID
	adj = ""
	// Retrieve the path after the resource ID
	// the index 0 of the splitted Path is "/"
	// the index 1 is the resource type (can differ from the resourceType)
	// the index 2 if exists is the resource id
	// the index 3 if exists is the adjective
	splittedPath := strings.Split(path, "/")
	pathLength := len(splittedPath)
	if pathLength >= 3 {
		idStr := strings.Replace(splittedPath[2], "/", "", -1)
		var err error
		id, err = primitive.ObjectIDFromHex(idStr)
		if err != nil {
			log.Printf("could not parse ID from URL : %s", err.Error())
		}
		if pathLength == 4 {
			adj = strings.Replace(splittedPath[3], "/", "", -1)
		}
	}
	return id, adj

}

// Abort just return the status code
func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

/*
* Method to handle Non Json Data (Basicaly Images)
 */
func (api *API) nonJSONrequestHandler(db *mongo.Database, resource NonJSONResource, resourceType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var content []byte
		var data interface{}
		var code int
		isData := false
		method := request.Method // Get HTTP Method (string)
		request.ParseForm()      // Populates request.Form
		values := request.Form
		id, adj := api.splitPath(request.URL.Path, resourceType)

		body := request.Body
		fmt.Printf("Received: %s with args : \n\t %+v\n", method, values)
		switch method {
		case http.MethodGet:
			code, data = resource.Get(db, values, id)
			if code == 200 {
				nonJSONResponse, ok := data.(NonJSONData)
				if ok {
					isData = true
					rw.Header().Set("Content-Type", nonJSONResponse.GetContentType())
					//rw.Header().Set("Content-Length", nonJSONResponse.GetContentLength())
					//rw.Header().Set("Accept-Ranges", "bytes")
					content = nonJSONResponse.GetContent()
				}
			}
		case http.MethodPost:
			code, data = resource.Post(db, values, request, id, adj)
		case http.MethodPut:
			code, data = resource.Put(db, values, body)
		case http.MethodDelete:
			code, data = resource.Delete(db, values, id)
		case http.MethodOptions:
			code = 200
			data = nil
		default:
			api.Abort(rw, 405)
		}
		if !isData {
			content, _ = json.Marshal(data)
			//if err != nil && method != GET {
			//	api.Abort(rw, 500)
			//}
		}

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		rw.WriteHeader(code)
		rw.Write(content)
	}
}
func (api *API) requestHandler(db *mongo.Database, resource data.Resource, resourceType string, isPrefixed bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		method := request.Method // Get HTTP Method (string)
		request.ParseForm()      // Populates request.Form
		values := request.Form
		splittedPath := strings.SplitAfter(request.URL.Path, "/")
		// Retrieve the path after the resource ID
		// the index 0 of the splitted Path is "/"
		// the index 1 is the resource type (can differ from the resourceType)(in case of a santard it's just standard)
		// the index 2 if exists is the resource id (in case of a standard it's the standard type)
		// the index 3 if exists is the adjective (incase of a standard it's the standard ID)
		id := primitive.NilObjectID
		adj := ""
		var err error
		pathLength := len(splittedPath)
		if isPrefixed {
			log.Print("We are a standard !")

			if pathLength >= 4 {
				id, err = primitive.ObjectIDFromHex(strings.Replace(splittedPath[3], "/", "", -1))
				if err != nil {
					log.Printf("Whow the given ID looks shitty : %s", err.Error())
					return
				}
				if pathLength == 5 {
					adj = strings.Replace(splittedPath[4], "/", "", -1)
				}
			}
		} else {
			if pathLength >= 3 {
				id, err = primitive.ObjectIDFromHex(strings.Replace(splittedPath[2], "/", "", -1))
				if err != nil {
					log.Printf("Whow the given ID looks shitty : %s", err.Error())
					return
				}
				if pathLength == 4 {
					adj = strings.Replace(splittedPath[3], "/", "", -1)
				}
			}

		}

		body := request.Body
		log.Printf("Received: %s with args : \n\t %+v\n", method, values)
		switch method {
		case http.MethodGet:
			code, data = resource.Get(db, values, id, adj)
		case http.MethodPost:
			code, data = resource.Post(db, values, request, id, adj)
		case http.MethodPut:
			code, data = resource.Put(db, values, body)
		case http.MethodDelete:
			code, data = resource.Delete(db, values, id)
		case http.MethodOptions:
			code = 200
			data = nil
		default:
			api.Abort(rw, 405)
		}
		content, err := json.Marshal(data)
		if err != nil {
			api.Abort(rw, 500)
		}
		rw.Header().Set("Content-Type", "text/json; charset=utf-8")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		rw.WriteHeader(code)
		rw.Write(content)
	}
}
func (api *API) returnStandardsLists() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var content []byte

		if request.Method != http.MethodGet {
			api.Abort(rw, 405)
		}

		content, err := json.Marshal(managedStandard)
		if err != nil {
			api.Abort(rw, 500)
		}
		rw.Header().Set("Content-Type", "text/json; charset=utf-8")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		rw.Header().Set("Access-Control-Allow-Methods", "GET")
		rw.WriteHeader(200)
		rw.Write(content)
	}
}

// AddResource add path to the http Handler
func (api *API) AddResource(db *mongo.Database, resource data.Resource, path string) {
	// Retrieve the Type Name of the Resource (Bike, Component etc ...)
	resourceType := strings.ToLower(reflect.TypeOf(resource).Elem().Name())
	subPath := ""
	if path != "" {
		subPath = fmt.Sprintf("%v/", path)
	} else {
		path = fmt.Sprintf("/%v", resourceType)
		subPath = fmt.Sprintf("/%v/", resourceType)

	}
	log.Printf("adding path %v", path)
	http.HandleFunc(path, api.requestHandler(db, resource, resourceType, false))
	http.HandleFunc(subPath, api.requestHandler(db, resource, resourceType, false))

}

// AddNonJSONResource will add a non json Handler
func (api *API) AddNonJSONResource(db *mongo.Database, resource NonJSONResource, path string) {
	resourceType := strings.ToLower(reflect.TypeOf(resource).Elem().Name())
	subPath := ""
	if path != "" {
		subPath = fmt.Sprintf("%v/", path)
	} else {
		path = fmt.Sprintf("/%v", resourceType)
		subPath = fmt.Sprintf("/%v/", resourceType)

	}
	log.Printf("adding path %v", path)
	http.HandleFunc(path, api.nonJSONrequestHandler(db, resource, resourceType))
	http.HandleFunc(subPath, api.nonJSONrequestHandler(db, resource, resourceType))

}

// Start starts the listener
func (api *API) Start(inetaddr string, port int) {
	portString := fmt.Sprintf("%s:%d", inetaddr, port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
