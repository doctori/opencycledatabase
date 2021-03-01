package api

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"gorm.io/gorm"
)

// Init will create the routes
func (api *API) Init(db *gorm.DB) {
	bike := new(data.Bike)
	component := new(data.Component)
	image := new(data.Image)
	// TODO : use Gorilla Mux ??
	api.AddResource(db, bike, "/bikes")
	api.AddResource(db, component, "/components")
	api.addStandard(db, &standards.BottomBracket{})
	api.addStandard(db, &standards.Crank{})
	api.addStandard(db, &standards.ChainRing{})
	api.AddResource(db, &data.Brand{}, "/brands")
	api.AddNonJSONResource(db, image, "")
	fmt.Printf("Listening To :8080 \n")
	api.Start("0.0.0.0", 8080)
}

// AddStandard add '/standards/%standardType%/ path to the http Handler
func (api *API) addStandard(db *gorm.DB, resource Resource) {
	// Retrieve the Type Name of the Resource (Bike, Component etc ...)
	resourceType := strings.ToLower(reflect.TypeOf(resource).Elem().Name())
	path := fmt.Sprintf("/standards/%v", resourceType)
	subPath := fmt.Sprintf("/standards/%v/", resourceType)
	log.Printf("adding path %v for resource %#v", path, resource)
	http.HandleFunc(path, api.requestHandler(db, resource, resourceType, true))
	http.HandleFunc(subPath, api.requestHandler(db, resource, resourceType, true))

}
