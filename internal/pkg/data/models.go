package data

import (
	"encoding/json"
	"errors"
	"io"

	// import png support
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"

	//"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const defaultPerPage int = 30

type Resource interface {
	Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{})
	Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{})
	Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *gorm.DB, values url.Values, id int) (int, interface{})
}

// Bike contains the defintion of the bike and all its attached components
type Bike struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	Name              string `gorm:"uniqueIndex:bike_uniqueness"`
	Brand             Brand
	BrandID           int    `json:"-" gorm:"uniqueIndex:bike_uniqueness"`
	Year              string `gorm:"uniqueIndex:bike_uniqueness"`
	Description       string
	Image             []Image              `gorm:"many2many:bike_images"`
	Components        []Component          `gorm:"many2many:bike_components;"`
	SupportedStandard []standards.Standard `gorm:"-"`
	PutNotSupported
}

// Should Return Errors !!
func (b Bike) save(db *gorm.DB) {
	// If we have a new record we create it
	if b.ID == 0 {
		oldb := new(Bike)
		b.Brand = b.Brand.save(db)
		db.Preload("Components").Preload("Brand").Where("name = ? AND brand_id = ? AND year = ?", b.Name, b.Brand.ID, b.Year).First(&oldb)
		log.Println(oldb)
		if oldb.Name == "" {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", b)
			log.Println("==========================================================================================")
			for i, component := range b.Components {
				component, _ := component.save(db)

				b.Components[i] = component
			}
			// We skip association since it can't say if an association already exists ...
			db.Set("gorm:save_associations", false).Create(&b)
			// We update our just created object in order to add it's associations ...
			db.Save(&b)
		} else {
			log.Println("Updating the record")
			for i, nc := range b.Components {

				nc, err := nc.save(db)
				if err != nil {
					log.Printf("Could not save : %v\n", nc)
				} else {
					log.Printf("Saved : %v\n", nc)
					b.Components[i] = nc
				}
			}
			db.Model(&oldb).Updates(&b)
			b = *oldb
		}
	} else {
		b.Brand = b.Brand.save(db)
		for i, component := range b.Components {
			component, _ := component.save(db)
			b.Components[i] = component
		}
		db.Save(&b)
	}
}

// GetAllBikes will return all the bikes with the pagination asked
func GetAllBikes(db *gorm.DB, page string, perPage string) interface{} {
	ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iperPage, err := strconv.Atoi(perPage)
	if err != nil {
		iperPage = defaultPerPage
	}
	var bikes []Bike
	//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
	db.Preload("Components").Preload("Brand").Order("name").Offset(ipage * iperPage).Limit(iperPage).Find(&bikes)
	return bikes
}

// Get Bike will return the request bike struct
func (Bike) Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{}) {
	page := values.Get("page")
	perPage := values.Get("per_page")
	/*if values.Get("name") == "" {
		return 200, GetAllBikes()
	}*/
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == 0 {
		return 200, GetAllBikes(db, page, perPage)
	}
	var bike Bike
	err := db.Preload("Components").Preload("Components.Standards").First(&bike, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Bike not found"
	}
	return 200, bike
}

// Delete the bike on DB given an bike ID -int-
func (Bike) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	if id != 0 {
		log.Print("Will Delete the ID : ")
		log.Println(id)
		err := db.Where("id = ?", id).Delete(&Bike{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 404, "Id Not Found"
		}
		return 200, "plop"
	}
	log.Println(id)
	return 404, "NOT FOUND"
}

// Post Handle Post request from http
func (Bike) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	var bike Bike
	if adj != "" {
		if id == 0 {
			return 500, "FUCK"
		}
		log.Printf("%#v", body)
		err := db.Preload("Brand").Preload("Components").Preload("Components.Standards").First(&bike, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 404, "Bike not found"
		}
	} else {
		decoder := json.NewDecoder(body)
		err := decoder.Decode(&bike)
		if err != nil {
			return 500, "Internal Error"
		}
		// Clean the unName Components
		components := bike.Components[:0]
		for _, component := range bike.Components {
			if component.Name != "" {
				components = append(components, component)
			}
		}
		bike.Components = components
		bike.save(db)
	}
	return 200, bike
}

// TODO : use it
func (Bike) addComponent(db *gorm.DB, bike Bike, component Component) {
	bike.Components = append(bike.Components, component)
	bike.save(db)
}

func (Bike) getCompatibleComponents(bike Bike) []Component {
	var compatibleComponents []Component
	// Get The Components, Get their Standards

	// Get All the Components supporting those standards
	// Return Thems
	return compatibleComponents
}
