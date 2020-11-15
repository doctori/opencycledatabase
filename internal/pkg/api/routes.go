package api

import (
	"fmt"

	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"gorm.io/gorm"
)

// Init will create the routes
func (api *API) Init(db *gorm.DB) {
	bike := new(data.Bike)
	component := new(data.Component)
	bbStandard := new(standards.BBStandard)
	image := new(data.Image)
	// TODO : use Gorilla Mux ??
	api.AddResource(db, bike, "/bikes")
	api.AddResource(db, component, "/components")
	api.AddResource(db, bbStandard, "")
	api.AddResource(db, &data.Brand{}, "/brands")
	api.AddNonJSONResource(db, image, "")
	fmt.Printf("Listening To :8080 \n")
	api.Start("0.0.0.0", 8080)
}
