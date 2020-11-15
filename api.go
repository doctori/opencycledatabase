package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Resource will define the generic API resource methods
type Resource interface {
	Get(values url.Values, id int) (int, interface{})
	Post(values url.Values, request *http.Request, id int, adj string) (int, interface{})
	Put(values url.Values, body io.ReadCloser) (int, interface{})
	Delete(values url.Values, id int) (int, interface{})
}

// NonJSONResource hold the files images and other objects
type NonJSONResource interface {
	Get(values url.Values, id int) (int, interface{})
	Post(values url.Values, request *http.Request, id int, adj string) (int, interface{})
	Put(values url.Values, body io.ReadCloser) (int, interface{})
	Delete(values url.Values, id int) (int, interface{})
}

// NonJSONData represent any content that is not JSon, mostly images
type NonJSONData interface {
	GetContentType() string
	GetContentLength() string
	GetContent() []byte
}
type (
	// GetNotSupported default response when you cannot Get the resource
	GetNotSupported struct{}
	// PostNotSupported default response when you cannot Post the resource
	PostNotSupported struct{}
	//PutNotSupported default response when you cannot Get the resource
	PutNotSupported struct{}
	// DeleteNotSupported default response when you cannot Get the resource
	DeleteNotSupported struct{}
)

// Get returns a 405
func (GetNotSupported) Get(values url.Values, id int) (int, interface{}) {
	return 405, ""
}

// Post returns a 405
func (PostNotSupported) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return 405, ""
}

// Put returns a 405
func (PutNotSupported) Put(values url.Values, body io.ReadCloser) (int, interface{}) {
	return 405, ""
}

// Delete returns a 405
func (DeleteNotSupported) Delete(values url.Values, id int) (int, interface{}) {
	return 405, ""
}

// API is the generic "API" model
type API struct{}

func (api *API) splitPath(path string, resourceType string) (id int, adj string) {
	id = 0
	adj = ""
	// Retrieve the path after the resource ID
	// the index 0 of the splitted Path is "/"
	// the index 1 is the resource type (can differ from the resourceType)
	// the index 2 if exists is the resource id
	// the index 3 if exiests is the adjective
	splittedPath := strings.Split(path, "/")
	pathLength := len(splittedPath)
	if pathLength >= 3 {
		id, _ = strconv.Atoi(strings.Replace(splittedPath[2], "/", "", -1))
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
func (api *API) nonJSONrequestHandler(resource NonJSONResource, resourceType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var content []byte
		var data interface{}
		var code int

		method := request.Method // Get HTTP Method (string)
		request.ParseForm()      // Populates request.Form
		values := request.Form
		id, adj := api.splitPath(request.URL.Path, resourceType)

		body := request.Body
		fmt.Printf("Received: %s with args : \n\t %+v\n", method, values)
		switch method {
		case http.MethodGet:
			var response interface{}
			code, response = resource.Get(values, id)
			if code == 200 {
				nonJSONResponse := response.(NonJSONData)
				rw.Header().Set("Content-Type", nonJSONResponse.GetContentType())
				//rw.Header().Set("Content-Length", nonJSONResponse.GetContentLength())
				//rw.Header().Set("Accept-Ranges", "bytes")
				content = nonJSONResponse.GetContent()
			}
		case http.MethodPost:
			code, data = resource.Post(values, request, id, adj)
		case http.MethodPut:
			code, data = resource.Put(values, body)
		case http.MethodDelete:
			code, data = resource.Delete(values, id)
		case http.MethodOptions:
			code = 200
			data = nil
		default:
			api.Abort(rw, 405)
		}
		if len(content) < 5 {
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
func (api *API) requestHandler(resource Resource, resourceType string) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		method := request.Method // Get HTTP Method (string)
		request.ParseForm()      // Populates request.Form
		values := request.Form
		splittedPath := strings.SplitAfter(request.URL.Path, "/")
		// Retrieve the path after the resource ID
		// the index 0 of the splitted Path is "/"
		// the index 1 is the resource type (can differ from the resourceType)
		// the index 2 if exists is the resource id
		// the index 3 if exists is the adjective
		id := 0
		adj := ""
		pathLength := len(splittedPath)
		if pathLength >= 3 {
			id, _ = strconv.Atoi(strings.Replace(splittedPath[2], "/", "", -1))
			if pathLength == 4 {
				adj = strings.Replace(splittedPath[3], "/", "", -1)
			}
		}

		body := request.Body
		fmt.Printf("Received: %s with args : \n\t %+v\n", method, values)
		switch method {
		case http.MethodGet:
			code, data = resource.Get(values, id)
		case http.MethodPost:
			code, data = resource.Post(values, request, id, adj)
		case http.MethodPut:
			code, data = resource.Put(values, body)
		case http.MethodDelete:
			code, data = resource.Delete(values, id)
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

// AddResource add path to the http Handler
func (api *API) AddResource(resource Resource, path string) {
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
	http.HandleFunc(path, api.requestHandler(resource, resourceType))
	http.HandleFunc(subPath, api.requestHandler(resource, resourceType))

}

// AddNonJSONResource will add a non json Handler
func (api *API) AddNonJSONResource(resource NonJSONResource, path string) {
	resourceType := strings.ToLower(reflect.TypeOf(resource).Elem().Name())
	subPath := ""
	if path != "" {
		subPath = fmt.Sprintf("%v/", path)
	} else {
		path = fmt.Sprintf("/%v", resourceType)
		subPath = fmt.Sprintf("/%v/", resourceType)

	}
	log.Printf("adding path %v", path)
	http.HandleFunc(path, api.nonJSONrequestHandler(resource, resourceType))
	http.HandleFunc(subPath, api.nonJSONrequestHandler(resource, resourceType))

}

// Start starts the listener
func (api *API) Start(inetaddr string, port int) {
	portString := fmt.Sprintf("%s:%d", inetaddr, port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
