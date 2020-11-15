package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	// import png support
	_ "image/png"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"github.com/nfnt/resize"

	//"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const defaultPerPage int = 30

// ComponentType is the defintion of the type of bike component
type ComponentType struct {
	gorm.Model
	Name        string
	Description string
}

// Brand hold the brand defintion
type Brand struct {
	gorm.Model
	Name         string
	Description  string
	Image        int
	CreationYear int
	EndYear      int
	Country      string
	PutNotSupported
	DeleteNotSupported
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

// Image is the description and the pointer to the image
type Image struct {
	gorm.Model
	Name          string `gorm:"uniqueIndex:image_uniqueness"`
	Path          string `gorm:"uniqueIndex:image_uniqueness"`
	Type          string
	ContentType   string
	ContentLength int64
	Content       []byte `sql:"-"`
	PutNotSupported
	DeleteNotSupported
}

// Get will return the asked Brand
func (Brand) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {

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
			return 500, "Internal Error"
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

// Get will return the image (content type is image/jpeg)
// TODO : make the content detection working
func (Image) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	var img Image
	if id == 0 {
		return 500, "Could not GET All Images"
	}
	err := db.First(&img, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Image Not Found"
	}
	file, err := os.Open(img.Path)
	log.Printf("Serving : %v", img.Path)
	defer file.Close()
	if err != nil {
		log.Println("Could Not Find the Image")
		return 404, "File Not Found"
	}
	// Detect Mime Type
	//buff := make([]byte, 512)
	//_, err = file.ReadAt(buff, 0)
	//img.ContentType = http.DetectContentType(buff)
	img.ContentType = "image/jpeg"

	decodedImg, _, err := image.Decode(file)
	file.Close()
	if err != nil {
		log.Panic(err)
	}
	resizedImg := resize.Resize(600, 0, decodedImg, resize.Lanczos3)
	var buffer bytes.Buffer
	jpeg.Encode(&buffer, resizedImg, nil)

	img.Content = buffer.Bytes()
	img.ContentLength = int64(buffer.Len())

	return 200, img

}
func (i Image) getContentType() string {
	return i.ContentType
}
func (i Image) getContentLength() string {
	return strconv.FormatInt(i.ContentLength, 10)
}

func (i Image) getContent() []byte {
	return i.Content
}

// Post saves the images to the "upload" directory (shouldn't it go to S3 ? )
func (Image) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	uploadFolder := "upload/"
	mediaType, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	checkErr(err, "Could Not Determine the Content Type")
	img := Image{}
	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("YOU ARE POSTING AN IMAGE")
		multipartReader := multipart.NewReader(request.Body, params["boundary"])
		// Should Build buffer and Write it at the end (loop thour the nextpart !)
		partReader, err := multipartReader.NextPart()
		fileName := partReader.FileName()
		filePath := uploadFolder + fileName
		file, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			return 500, err
		}
		checkErr(err, "Could Not Get the Next Part")
		if _, err := io.Copy(file, partReader); err != nil {
			log.Fatal(err)
			return 500, err
		}
		// Let's record this image and return it to our client
		img = Image{Name: fileName, Path: filePath}
		recordedImg := img.save(db)
		log.Printf("Posted : %#v", img)
		return 200, recordedImg
	}

	return 200, img
}

func (i Image) save(db *gorm.DB) Image {
	if i.ID == 0 {
		oldi := new(Image)
		db.Where("name = ? and path = ? ", i.Name, i.Path).First(&oldi)
		if oldi.Name == "" {
			log.Println("Recording the New Image")
			db.Create(&i)
		} else {
			log.Println("Updating The Image Record")
			db.Model(&oldi).Updates(&i)
			i = *oldi
			log.Printf("Saving Image %#v", i)
		}
	} else {
		db.Create(&i)
	}
	return i
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
func (Bike) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
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
