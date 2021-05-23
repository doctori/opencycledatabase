package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	managedStandard []data.Resource
)

// Init will create the routes
func (api *API) Init(db *mongo.Database, conf *config.Config) {
	bike := new(data.Bike)
	component := new(data.Component)
	image := new(data.Image)

	// static content
	api.addStaticDir("./upload")
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
	api.addStandard(db, standards.NewFrame())
	api.addStandard(db, standards.NewSeatTube())
	api.AddResource(db, &data.Brand{}, "/brands")
	http.HandleFunc("/standards", api.returnStandardsLists())
	api.AddNonJSONResource(db, image, "/images")

	// Once all standards have been routed, migrate their model
	log.Infof("Listening To %s:%d \n", conf.API.BindIP, conf.API.BindPort)
	api.Start(conf.API.BindIP, conf.API.BindPort)
}

type StatusRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *StatusRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	w.ResponseWriter.WriteHeader(status)
}
func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srw := &StatusRespWr{ResponseWriter: w}
		h.ServeHTTP(srw, r)
		if srw.status >= 400 { // 400+ codes are the error codes
			log.Printf("Error status code: %d when serving path: %s",
				srw.status, r.RequestURI)
		}
	}
}
func (api *API) addStaticDir(directory string) {
	// static content server
	fs := http.FileServer(http.Dir(directory))

	directory = strings.TrimLeft(directory, ". ")
	log.Debugf("Adding static directory %s", directory)
	http.Handle("/upload/", wrapHandler(http.StripPrefix(directory, fs)))

}

// AddStandard add '/standards/%standardType%/ path to the http Handler
func (api *API) addStandard(db *mongo.Database, resource data.Resource) {
	std := resource.(standards.StandardInt)
	// Retrieve the Type Name of the Resource (Bike, Component etc ...)
	resourceType := std.GetType()
	//stdType.Save(db)
	managedStandard = append(managedStandard, resource)
	path := fmt.Sprintf("/standards/%v", strings.ToLower(resourceType))
	subPath := fmt.Sprintf("/standards/%v/", strings.ToLower(resourceType))
	componentsPath := fmt.Sprintf("/standards/%v/components", strings.ToLower(resourceType))
	log.Debugf("adding path [%s] for resource %s", path, resourceType)
	http.HandleFunc(path, api.requestHandler(db, resource, resourceType, true))
	http.HandleFunc(subPath, api.requestHandler(db, resource, resourceType, true))
	http.HandleFunc(componentsPath, api.requestHandler(db, resource, resourceType, true))

}
