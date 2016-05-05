package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	//"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const PER_PAGE int = 30

type ComponentType struct {
	gorm.Model
	Name        string
	Description string
}
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
type Bike struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	Name              string
	Brand             Brand
	BrandID           int `json:"-"`
	Year              string
	Description       string
	Image             int
	Components        []Component `gorm:"many2many:bike_components;"`
	SupportedStandard []Standard  `sql:"-"`
	PutNotSupported
}
type Component struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	Name        string
	Brand       Brand
	BrandID     int `json:"-"`
	Type        ComponentType
	TypeID      int `json:"-"`
	Description string
	Standards   []Standard `gorm:"many2many:component_standards"`
	Year        string
	PutNotSupported
	DeleteNotSupported
}
type Image struct {
	gorm.Model
	Name          string
	Path          string
	Type          string
	ContentType   string
	ContentLength int64
	Content       []byte `sql:"-"`
	PutNotSupported
	DeleteNotSupported
}
type Standard struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	Name        string
	Country     string
	Code        string
	Type        string
	Description string
}

var db = &gorm.DB{}

func initDB(config Config) *gorm.DB {
	connectionString := fmt.Sprintf(
		"user=%s password='%s' host=%s dbname=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.DBname)
	log.Printf("Connecting to %s", connectionString)
	db, err := gorm.Open("postgres", connectionString)
	checkErr(err, "Postgres Opening Failed")
	// Debug Mode
	db.LogMode(true)
	db.CreateTable(&Image{}, &ComponentType{}, &Brand{}, &Standard{}, &Component{}, &Bike{})
	db.Model(&Bike{}).AddUniqueIndex("bike_uniqueness", "name,  year")
	db.Model(&Component{}).AddUniqueIndex("component_uniqueness", "name, year")
	db.Model(&Standard{}).AddUniqueIndex("standard_uniqueness", "name, code, type")
	db.Model(&Image{}).AddUniqueIndex("image_uniqueness", "name", "path")
	db.AutoMigrate(&Bike{}, &Component{}, &Standard{}, &Image{}, &Brand{}, &ComponentType{})
	checkErr(err, "Create tables failed")

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}

func (Brand) Get(values url.Values, id int) (int, interface{}) {

	if id == 0 {
		var brands []Brand
		//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
		db.Find(&brands)
		return 200, brands
	}
	var brand Brand
	err := db.First(&brand, id).RecordNotFound()
	if err {
		return 404, "Brand not found"
	}
	return 200, brand
}

func (Brand) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	var brand Brand
	if adj != "" {
		if id == 0 {
			return 500, "That Shouldn't have appended"
		}

		err := db.First(&brand, id).RecordNotFound()
		if err {
			return 404, "Brand not found"
		}
	} else {
		decoder := json.NewDecoder(body)
		err := decoder.Decode(&brand)
		if err != nil {
			panic(err)
			return 500, "Internal Error"
		}
		log.Println(brand)
		brand = brand.save()
	}

	return 200, brand
}

func (b Brand) save() Brand {
	if db.NewRecord(b) {
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

func (ct ComponentType) save() ComponentType {
	if db.NewRecord(ct) {
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

func (Image) Get(values url.Values, id int) (int, interface{}) {
	var img Image
	if id == 0 {
		return 500, "Could not GET All Images"
	}
	recordNotFound := db.First(&img, id).RecordNotFound()
	if recordNotFound {
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
func (i Image) GetContentType() string {
	return i.ContentType
}
func (i Image) GetContentLength() string {
	return strconv.FormatInt(i.ContentLength, 10)
}
func (i Image) GetContent() []byte {
	return i.Content
}
func (Image) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
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
		recordedImg := img.save()
		log.Printf("Posted : %#v", img)
		return 200, recordedImg
	} else {
		return 500, "Nhope"
	}

	return 200, img
}

func (i Image) save() Image {
	if db.NewRecord(i) {
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
func (b Bike) save() {
	// If we have a new record we create it
	if db.NewRecord(b) {
		oldb := new(Bike)
		b.Brand = b.Brand.save()
		db.Preload("Components").Preload("Brand").Where("name = ? AND brand_id = ? AND year = ?", b.Name, b.Brand.ID, b.Year).First(&oldb)
		log.Println(oldb)
		if oldb.Name == "" {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", b)
			log.Println("==========================================================================================")
			for i, component := range b.Components {
				_, component := component.save()

				b.Components[i] = component
			}
			// We skip association since it can't say if an association already exists ...
			db.Set("gorm:save_associations", false).Create(&b)
			// We update our just created object in order to add it's associations ...
			db.Save(&b)
		} else {
			log.Println("Updating the record")
			for i, nc := range b.Components {

				err, nc := nc.save()
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
		b.Brand = b.Brand.save()
		for i, component := range b.Components {
			_, component := component.save()
			b.Components[i] = component
		}
		db.Save(&b)
	}
}

func GetAllBikes(page string, per_page string) interface{} {
	ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iper_page, err := strconv.Atoi(per_page)
	if err != nil {
		iper_page = PER_PAGE
	}
	var bikes []Bike
	//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
	db.Preload("Components").Preload("Brand").Order("name").Offset(ipage * iper_page).Limit(iper_page).Find(&bikes)
	return bikes
}

func (Bike) Get(values url.Values, id int) (int, interface{}) {
	page := values.Get("page")
	per_page := values.Get("per_page")
	/*if values.Get("name") == "" {
		return 200, GetAllBikes()
	}*/
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == 0 {
		return 200, GetAllBikes(page, per_page)
	}
	var bike Bike
	err := db.Preload("Components").Preload("Components.Standards").First(&bike, id).RecordNotFound()
	if err {
		return 404, "Bike not found"
	}
	return 200, bike
}

func (Bike) Delete(values url.Values, id int) (int, interface{}) {
	if id != 0 {
		log.Print("Will Delete the ID : ")
		log.Println(id)
		err := db.Where("id = ?", id).Delete(&Bike{}).RecordNotFound()
		if err {
			return 404, "Id Not Found"
		}
		return 200, "plop"
	} else {
		log.Println(id)
		return 404, "NOT FOUND"
	}
}

func (Bike) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	var bike Bike
	if adj != "" {
		if id == 0 {
			return 500, "FUCK"
		}
		log.Printf("%#v", body)
		err := db.Preload("Brand").Preload("Components").Preload("Components.Standards").First(&bike, id).RecordNotFound()
		if err {
			return 404, "Bike not found"
		}
	} else {
		decoder := json.NewDecoder(body)
		err := decoder.Decode(&bike)
		if err != nil {
			panic(err)
			return 500, "Internal Error"
		}
		// Clean the unName Components
		components := bike.Components[:0]
		for _, component := range bike.Components {
			if component.Name != "" {
				components = append(components, component)
			}
		}
		bike.save()
	}
	return 200, bike
}

func (Bike) addComponent(bike Bike, component Component) {
	bike.Components = append(bike.Components, component)
	bike.save()
}

func (Bike) getCompatibleComponents(bike Bike) []Component {
	var compatibleComponents []Component
	// Get The Components, Get their Standards

	// Get All the Components supporting those standards
	// Return Thems
	return compatibleComponents
}

func (Standard) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	var standard Standard
	err := decoder.Decode(&standard)
	if err != nil {
		panic("shiit")
	}
	log.Println(standard)
	err, standard = standard.save()
	if err != nil {
		return 500, "Could Not Save the Standard"
	}
	return 200, standard
}

func (Standard) Put(values url.Values, body io.ReadCloser) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	var standard Standard
	err := decoder.Decode(&standard)
	if err != nil {
		panic("shiit")
	}
	log.Println(standard)
	err, standard = standard.save()
	if err != nil {
		return 500, "Could Not Save the Standard"
	}
	return 200, standard
}
func (Standard) Get(values url.Values, id int) (int, interface{}) {
	filename := "db/standard_" + values.Get("name") + ".json"
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return 404, "404 Standard Not Found"
	}
	var s Standard
	if json.Unmarshal(jsonBlob, &s) != nil {
		return 500, ""
	}

	return 200, s
}
func (Standard) Delete(values url.Values, id int) (int, interface{}) {
	filename := "db/standard_" + values.Get("name") + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return 404, "404 Standard Not Found"
	}
	if os.Remove(filename) != nil {
		return 500, "Could not Delete the Standard"
	}
	return 200, ""
}
func (s Standard) save() (error, Standard) {
	filename := "db/standard_" + s.Name + ".json"
	//if (s.ID == ""){
	//	s.ID = uuid.NewV4().String()
	//}
	jsonBlob, err := json.Marshal(s)
	if err != nil {
		return err, s
	}
	err = ioutil.WriteFile(filename, jsonBlob, 0600)
	return err, s
}

func (Component) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	fmt.Printf("Received args : \n\t %+v\n", body)
	decoder := json.NewDecoder(body)
	var component Component
	err := decoder.Decode(&component)
	if err != nil {
		log.Printf("Could not decode : %+v\nbecause %v\n", body, err)
		return 500, "Could Not Decode the Data"
		panic("shiit")
	}
	log.Println(component)
	err, component = component.save()
	if err != nil {
		return 500, "Could Not Save the Component"
	}
	return 200, component
}

func (Component) Get(values url.Values, id int) (int, interface{}) {
	page := values.Get("page")
	per_page := values.Get("per_page")
	if values.Get("name") == "" {
		c := new(Component)
		return 200, c.getAll(page, per_page)
	}

	log.Println(values.Get("name"))
	var c Component
	err := db.Preload("Standards").Preload("Type").Preload("Brand").Find(&c, "name= ? ", values.Get("name")).RecordNotFound()
	if err {
		return 404, "Component not found"
	}
	return 200, c
}
func (Component) getAll(page string, per_page string) []Component {
	ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iper_page, err := strconv.Atoi(per_page)
	if err != nil {
		iper_page = PER_PAGE
	}

	var components []Component
	//TODO : return LINKS Header with the next page and previous page
	db.Preload("Standards").Preload("Type").Preload("Brand").Offset(ipage * iper_page).Limit(iper_page).Find(&components)
	return components
}

// Return Components Compatibles with a standards
func (Component) getCompatible(s Standard) []Component {
	var compatibleComponents []Component
	// Crado Way
	c := new(Component)
	components := c.getAll("0", "30")
	for _, component := range components {
		for _, standard := range component.Standards {
			if standard == s {
				compatibleComponents = append(compatibleComponents, component)
			}
		}
	}
	return compatibleComponents
}

func (c Component) save() (error, Component) {
	if db.NewRecord(c) {
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
			c.Brand = c.Brand.save()
			c.Type = c.Type.save()
			db.Model(&oldc).Updates(&c)
			c = *oldc
		}
	} else {
		// Maybe Save nested object independantly ?
		c.Type = c.Type.save()
		c.Brand = c.Brand.save()
		db.Save(&c)
	}
	return db.Error, c
}
