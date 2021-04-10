package api

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"gorm.io/gorm"
)

var (
	managedStandard []Resource
)

// Init will create the routes
func (api *API) Init(db *gorm.DB, conf *config.Config) {
	bike := new(data.Bike)
	component := new(data.Component)
	image := new(data.Image)

	// TODO : use Gorilla Mux ??
	api.AddResource(db, bike, "/bikes")
	api.AddResource(db, component, "/components")
	api.addStandard(db, standards.NewBottomBracket())
	api.addStandard(db, standards.NewCrank())
	api.addStandard(db, standards.NewChainRing())
	api.addStandard(db, standards.NewFrontDerailleur())
	api.addStandard(db, standards.NewHeadset())
	api.addStandard(db, standards.NewHub())
	api.addStandard(db, standards.NewRearDerailleur())
	api.addStandard(db, standards.NewWheel())
	api.addStandard(db, standards.NewSpoke())
	api.addStandard(db, standards.NewThread())
	api.AddResource(db, &data.Brand{}, "/brands")
	http.HandleFunc("/standards", api.returnStandardsLists())
	api.AddNonJSONResource(db, image, "/images")
	fmt.Printf("Listening To %s:%d \n", conf.API.BindIP, conf.API.BindPort)
	api.Start(conf.API.BindIP, conf.API.BindPort)
}

// AddStandard add '/standards/%standardType%/ path to the http Handler
func (api *API) addStandard(db *gorm.DB, resource Resource) {

	// Retrieve the Type Name of the Resource (Bike, Component etc ...)
	resourceType := strings.ToLower(reflect.TypeOf(resource).Elem().Name())
	managedStandard = append(managedStandard, resource)
	path := fmt.Sprintf("/standards/%v", resourceType)
	subPath := fmt.Sprintf("/standards/%v/", resourceType)
	log.Printf("adding path %v for resource %#v", path, resource)
	http.HandleFunc(path, api.requestHandler(db, resource, resourceType, true))
	http.HandleFunc(subPath, api.requestHandler(db, resource, resourceType, true))

}
