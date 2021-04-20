package data

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// Brand hold the brand defintion
type Brand struct {
	gorm.Model
	Name         string
	Description  string
	Image        int
	CreationYear int
	EndYear      int
	Country      string
	WikiHref     string
	Href         string
	PutNotSupported
	DeleteNotSupported
}

// Get will return the asked Brand
func (Brand) Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{}) {

	if id == 0 {
		var brands []Brand
		//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
		db.Find(&brands)
		return 200, brands
	}
	var brand Brand
	err := db.First(&brand, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Brand not found"
	}
	return 200, brand
}

// Post will save the brand
func (Brand) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	var brand Brand
	if adj != "" {
		if id == 0 {
			return 500, "That Shouldn't have appended"
		}

		err := db.First(&brand, id).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 404, "Brand not found"
		}
	} else {
		decoder := json.NewDecoder(body)
		err := decoder.Decode(&brand)
		if err != nil {
			log.Println(err)
			return 500, err.Error()
		}
		log.Println(brand)
		brand = brand.save(db)
	}

	return 200, brand
}

func (b Brand) save(db *gorm.DB) Brand {
	if b.ID == 0 {
		oldb := new(Brand)
		db.Where("name = ?", b.Name).First(&oldb)
		if oldb.Name == "" {
			log.Println("Recording the New Brand")
			db.Create(&b)
		} else {
			log.Println("Updating The Existing  Brand")
			db.Model(&oldb).Updates(&b)
			b = *oldb
			log.Printf("Saving Brand %#v", b)
		}
	} else {
		log.Printf("Creating Brand %#v", b)
		db.Save(&b)
	}
	return b
}
