package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"gorm.io/gorm"
)

// Component : Generic struct to regroup most common properties
type Component struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	Name        string `gorm:"uniqueIndex:component_uniqueness"`
	Brand       Brand
	BrandID     int `json:"-" gorm:"uniqueIndex:component_uniqueness"`
	Type        ComponentType
	TypeID      int `json:"-"`
	Description string
	Standards   []standards.Standard `gorm:"many2many:component_standards"`
	Images      []Image              `gorm:"many2many:component_images"`
	Year        string               `gorm:"uniqueIndex:component_uniqueness"`
	PutNotSupported
	DeleteNotSupported
}

// ComponentInt : the standard COmponent interface that needs to complies to in order to be a component
type ComponentInt interface {
	GetName() string
	GetBrand() Brand
	GetType() ComponentType
	GetDescription() string
	GetStandards() []standards.StandardInt
	GetImages() []Image
}

func (ct ComponentType) save(db *gorm.DB) ComponentType {
	if ct.ID == 0 {
		oldct := new(ComponentType)
		db.Where("name = ?", oldct.Name).First(&oldct)
		if oldct.Name == "" {
			log.Println("Recording the New Component type")
			db.Create(&ct)
		} else {
			log.Println("Updating The Existing Component Type")
			db.Model(&oldct).Updates(&ct)
			ct = *oldct
			log.Printf("Saving Component type %#v", ct)
		}
	} else {
		db.Save(&ct)
	}
	return ct
}

// Post will save the component in database
func (Component) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	fmt.Printf("Received args : \n\t %+v\n", body)
	decoder := json.NewDecoder(body)
	var component Component
	err := decoder.Decode(&component)
	if err != nil {
		log.Printf("Could not decode : %+v\nbecause %v\n", body, err)
		return 500, "Could Not Decode the Data"
	}
	log.Println(component)
	component, err = component.save(db)
	if err != nil {
		return 500, "Could Not Save the Component"
	}
	return 200, component
}

// Get return a generic Component
func (Component) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {

	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		page = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	perPage, err := strconv.Atoi(values.Get("per_page"))
	if err != nil {
		perPage = defaultPerPage
	}
	var c Component
	if values.Get("name") == "" && values.Get("search") == "" {
		return 200, c.getAll(db, page, perPage)
	} else if values.Get("search") != "" {
		return 200, c.search(db, page, perPage, values.Get("search"))
	}

	log.Println(values.Get("name"))
	err = db.Preload("Standards").
		Preload("Type").
		Preload("Brand").
		Find(&c, "name= ? ", values.Get("name")).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Component not found"
	}
	return 200, c
}

func (Component) search(db *gorm.DB, page int, perPage int, filter string) (components []Component) {
	filter = fmt.Sprintf("%%%s%%", filter)
	db.Preload("Standards").
		Preload("Type").
		Preload("Brand").
		Where("name LIKE ?", filter).
		Find(&components).
		Offset(page * perPage).
		Limit(perPage)
	return
}
func (Component) getAll(db *gorm.DB, page int, perPage int) (components []Component) {

	//TODO : return LINKS Header with the next page and previous page
	db.Preload("Standards").Preload("Type").Preload("Brand").Offset(page * perPage).Limit(perPage).Find(&components)
	return components
}

func (c Component) save(db *gorm.DB) (Component, error) {
	if c.ID == 0 {
		oldc := new(Component)
		db.Where(c.Brand).First(&c.Brand)
		db.Where(c.Type).First(&c.Type)

		if c.Brand.ID != 0 {
			log.Printf("Looking for : name = %v AND brand_id = %v AND type_id = %v AND year = %v", c.Name, c.Brand.ID, c.Type.ID, c.Year)
			db.Preload("Standards").Preload("Type").Preload("Brand").Where("name = ? AND brand_id =  ? AND year = ?", c.Name, c.Brand.ID, c.Year).First(&oldc)
		} else if c.Type.ID != 0 {
			log.Printf("Looking for : name = %v AND brand_id = %v AND type_id = %v AND year = %v", c.Name, c.Brand.ID, c.Type.ID, c.Year)
			db.Preload("Standards").Preload("Type", "id = ?", c.Type.ID).Where("name = ? AND year = ?", c.Name, c.Year).First(&oldc)
		} else {
			log.Printf("Looking for : name = %v AND brand_id = %v AND type_id = %v AND year = %v", c.Name, c.Brand.ID, c.Type.ID, c.Year)
			db.Preload("Standards").Preload("Type").Preload("Brand").Where("name = ? AND year = ?", c.Name, c.Year).First(&oldc)
		}
		log.Println(oldc)
		if oldc.Name == "" {
			log.Println("Creating the Component !")
			db.Create(&c)
		} else {
			log.Println("Updating the Component ...")
			// Let's check that for our Standard ...
			for j, ns := range c.Standards {
				if ns.ID == 0 {
					for _, s := range oldc.Standards {
						if s.Name == ns.Name && s.Code == ns.Code && s.Type == ns.Type {
							c.Standards[j].ID = s.ID
						}
					}
				}
			}
			c.Brand = c.Brand.save(db)
			c.Type = c.Type.save(db)
			db.Model(&oldc).Updates(&c)
			c = *oldc
		}
	} else {
		// Maybe Save nested object independantly ?
		c.Type = c.Type.save(db)
		c.Brand = c.Brand.save(db)
		db.Save(&c)
	}
	return c, db.Error
}
