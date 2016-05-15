package models

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
	"reflect"
	"strconv"
	"strings"
	//"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const PER_PAGE int = 30

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(values url.Values, id int) (int, interface{}) {
	return 405, ""
}

func (PostNotSupported) Post(values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return 405, ""
}

func (PutNotSupported) Put(values url.Values, body io.ReadCloser) (int, interface{}) {
	return 405, ""
}

func (DeleteNotSupported) Delete(values url.Values, id int) (int, interface{}) {
	return 405, ""
}

type DBConfig struct {
	Host     string `json:host`
	Username string `json:username`
	Password string `json:password`
	Port     int    `json:port`
	DBname   string `json:dbname`
}
type Config struct {
	DB DBConfig `json:db`
}

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
	Images      []Image    `gorm:"many2many:component_images"`
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

var (
	DB *gorm.DB
)

func init() {
	// Load Application Configuration info our Config Struct
	file, err := os.Open("conf/mainConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	log.Printf("%#v", config)
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.OpenFile("model_logs.log", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Panicln("Could Not Open the log File !!")
		log.Println(err)
		os.Exit(4)
	}
	log.SetOutput(logFile)
	DB, err = connectDB(config.DB)
	checkErr(err, "Postgres Opening Failed")
	// Debug Mode
	DB.LogMode(true)

	DB.CreateTable(&Image{}, &ComponentType{}, &Brand{}, &Standard{}, &Component{}, &Bike{})
	DB.Model(&Bike{}).AddUniqueIndex("bike_uniqueness", "name,  year, brand_id")
	DB.Model(&Component{}).AddUniqueIndex("component_uniqueness", "name, year, brand_id")
	DB.Model(&Standard{}).AddUniqueIndex("standard_uniqueness", "name, code, type")
	DB.Model(&Image{}).AddUniqueIndex("image_uniqueness", "name", "path")
	DB.AutoMigrate(&Bike{}, &Component{}, &Standard{}, &Image{}, &Brand{}, &ComponentType{})
	checkErr(err, "Create tables failed")

}
func connectDB(config DBConfig) (db *gorm.DB, err error) {
	log.Printf("Connecting to %s", config.Host)
	db, err = gorm.Open("postgres", fmt.Sprintf(
		"user=%s password='%s' host=%s dbname=%s",
		config.Username,
		config.Password,
		config.Host,
		config.DBname))
	return
}
func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}

func (Brand) Get(values url.Values, id int) (int, interface{}) {
	if id == 0 {
		var brands []Brand
		//DB.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
		DB.Find(&brands)
		return 200, brands
	}
	var brand Brand
	err := DB.First(&brand, id).RecordNotFound()
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

		err := DB.First(&brand, id).RecordNotFound()
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
		brand.Save()
	}

	return 200, brand
}

// Returns True if the left Brand Equals the Right Brand
func (lb *Brand) Equals(rb *Brand) bool {
	if lb.Name == rb.Name && lb.Description == rb.Description && lb.ID == rb.ID {
		return true
	} else {
		return false
	}
}
func (b *Brand) Save() {
	// if b has no ID
	if DB.NewRecord(b) {
		oldb := new(Brand)
		// Let's find if it exists
		DB.Where("name = ?", b.Name).First(oldb)
		if oldb.ID == 0 && oldb.Name == "" {
			log.Println("Recording the New Brand")
			DB.Create(b)
		} else {
			log.Println("Updating The Existing  Brand")
			/*DB.Model(oldb).Updates(b)
			b = oldb*/
			DB.Model(b).Updates(oldb)
			log.Printf("Saving Brand %#v\n", b)
		}
	} else {
		log.Printf("Updating Brand %#v\n", b)
		DB.Save(b)
	}
	log.Printf("Brand %#v has been SAVED\n", b)
}

func (ct *ComponentType) Save() {
	if DB.NewRecord(ct) {
		oldct := new(ComponentType)
		DB.Where("name = ?", oldct.Name).First(&oldct)
		if oldct.Name == "" {
			log.Println("Recording the New Component type")
			DB.Create(ct)
		} else {
			log.Println("Updating The Existing Component Type")
			DB.Model(&oldct).Updates(ct)
			ct = oldct
			log.Printf("Saving Component type %#v", *ct)
		}
	} else {
		DB.Save(ct)
	}
}

// Should Return Errors !!
func (b *Bike) Save() {
	// If we have a new record we create it
	if DB.NewRecord(b) {
		oldb := new(Bike)
		b.Brand.Save()
		log.Printf(" >>>>>>BRAND IS : %#v\n", b.Brand)
		log.Printf("Looking For : bike WHERE name = %v AND brand_id = %v AND year = %v", b.Name, b.Brand.ID, b.Year)
		DB.Preload("Brand").Where("name = ? AND brand_id = ? AND year = ?", b.Name, b.Brand.ID, b.Year).First(oldb)
		log.Println(oldb)
		if oldb.Name == "" {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", b)
			log.Println("==========================================================================================")
			for i := range b.Components {
				b.Components[i].Save()
			}
			// We skip association since it can't say if an association already exists ...
			DB.Set("gorm:save_associations", false).Create(&b)
			// We update our just created object in order to add it's associations ...
			DB.Save(b)
		} else {
			log.Println("Updating the record")
			for i := range b.Components {

				err := b.Components[i].Save()
				if err != nil {
					log.Printf("Could not save : %v\n", b.Components[i])
				} else {
					log.Printf("Saved : %v\n", b.Components[i])
				}
			}
			// Reflect our bikes in order to loop throught its properties
			oldbR := reflect.ValueOf(oldb).Elem()
			bR := reflect.ValueOf(b).Elem()
			for i := 0; i < bR.NumField(); i++ {
				// We retrieve the current element from our new Bike
				currElem := bR.Field(i).Interface()
				// If empty we put the old values
				if (currElem == nil) || (currElem == 0 || (currElem == "") || (currElem == false)) {
					bR.Field(i).Set(reflect.Value(oldbR.Field(i)))
				}
			}
			DB.Save(b)
		}
	} else {
		b.Brand.Save()
		for i := range b.Components {
			b.Components[i].Save()
		}
		DB.Save(b)
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
	DB.Preload("Components").Preload("Brand").Preload("Components.Type").Preload("Components.Brand").Order("name").Offset(ipage * iper_page).Limit(iper_page).Find(&bikes)
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
	err := DB.Preload("Components").Preload("Components.Standards").First(&bike, id).RecordNotFound()
	if err {
		return 404, "Bike not found"
	}
	return 200, bike
}

func (Bike) Delete(values url.Values, id int) (int, interface{}) {
	if id != 0 {
		log.Print("Will Delete the ID : ")
		log.Println(id)
		err := DB.Where("id = ?", id).Delete(&Bike{}).RecordNotFound()
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
		err := DB.Preload("Brand").Preload("Components").Preload("Components.Standards").First(&bike, id).RecordNotFound()
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
			if len(component.Name) > 3 {
				components = append(components, component)
			}
		}
		bike.Components = components
		bike.Save()
	}
	return 200, bike
}

func (Bike) addComponent(bike Bike, component Component) {
	bike.Components = append(bike.Components, component)
	bike.Save()
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
	errs := standard.Save()
	if errs != nil {
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
	errs := standard.Save()
	if errs != nil {
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
func (s *Standard) Save() []error {
	var err []error
	// If our Standard doesn't have an ID
	if DB.NewRecord(s) {
		olds := new(Standard)
		DB.Where("name = ?", s.Name).First(&olds)
		if olds.ID == 0 {
			log.Println("Recording the New Standard")
			err = DB.Create(s).GetErrors()
		} else {
			log.Println("Updating The Image Record")
			err = DB.Model(&olds).Updates(s).GetErrors()
			s = olds
			log.Printf("Saved Standard %#v\n", s)
		}
	} else {
		err = DB.Save(s).GetErrors()
	}
	return err
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
	err = component.Save()
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
	err := DB.Preload("Standards").Preload("Type").Preload("Brand").Find(&c, "name= ? ", values.Get("name")).RecordNotFound()
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
	DB.Preload("Standards").Preload("Type").Preload("Brand").Offset(ipage * iper_page).Limit(iper_page).Find(&components)
	return components
}

// Return Components Compatibles with a standards
func (Component) getCompatible(s Standard) []Component {
	var compatibleComponents []Component
	// Crado Way
	c := new(Component)
	// SRSLY ?
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

func (c *Component) Save() error {
	// the component has no ID
	if DB.NewRecord(c) {
		oldc := new(Component)

		// Let's find the brand and Type
		DB.Where("name = ? ", c.Brand.Name).First(&(c.Brand))
		if c.Type.Name != "" {
			DB.Where("name = ? ", c.Type.Name).First(&(c.Type))
		}

		if c.Brand.ID != 0 {
			// We find the brand end the type.
			if c.Type.ID != 0 {
				log.Printf("Looking for : name = %v AND brand_id = %v AND type_id = %v AND year = %v", c.Name, c.Brand.ID, c.Type.ID, c.Year)
				DB.Preload("Standards").Preload("Type", "id = ?", c.Type.ID).Where("name = ? AND year = ?", c.Name, c.Year).First(&oldc)
			} else {
				log.Printf("Looking for : name = %v AND brand_id = %v AND type_id = %v AND year = %v", c.Name, c.Brand.ID, c.Type.ID, c.Year)
				DB.Preload("Standards").Preload("Type").Preload("Brand").Where("name = ? AND brand_id =  ? AND type_id = ? AND year = ?", c.Name, c.Brand.ID, c.Type.ID, c.Year).First(&oldc)
			}
		} else if c.Type.ID != 0 {
			log.Printf("Looking for : name = %v AND type_id = %v AND year = %v", c.Name, c.Type.ID, c.Year)
			DB.Preload("Standards").Preload("Type").Where("name = ? AND type_id = ? AND year = ?", c.Name, c.Type.ID, c.Year).First(&oldc)

		} else {
			log.Printf("Looking for : name = %v AND brand = %v AND type = %v", c.Name, c.Brand.Name, c.Type.Name, c.Year)
			DB.Preload("Standards").Preload("Type", "name = ?", c.Type.Name).Preload("Brand", "name = ?", c.Brand.Name).Where("name = ? AND year = ?", c.Name, c.Year).First(&oldc)
		}
		log.Println(oldc)
		if oldc.Name == "" {
			log.Println("Creating the Component !")
			DB.Create(c)
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
			// Save the Images (Db Style)
			for i := range c.Images {
				c.Images[i].Save()
			}

			log.Printf(" >>>>>>>>>>>>> The Component Brand %#v Is Going to Be Saved been Saved", c.Brand)
			c.Brand.Save()
			log.Printf(" >>>>>>>>>>>>> The Component Brand %#v Has been Saved", c.Brand)
			c.Type.Save()
			DB.Model(&oldc).Updates(c)
		}
	} else {
		// Maybe Save nested object independantly ?
		c.Type.Save()
		c.Brand.Save()
		DB.Save(c)
	}
	return DB.Error
}

func (Image) Get(values url.Values, id int) (int, interface{}) {
	var img Image
	if id == 0 {
		return 500, "Could not GET All Images"
	}
	recordNotFound := DB.First(&img, id).RecordNotFound()
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
		img.Save()
		log.Printf("Posted : %#v", img)
		return 200, img
	} else {
		return 500, "Nhope"
	}

	return 200, img
}

func (img *Image) Save() {

	if DB.NewRecord(img) {
		oldi := new(Image)
		DB.Where("name = ?", img.Name).First(oldi)
		if oldi.ID == 0 && oldi.Name == "" {
			log.Println("Recording the New Image")
			DB.Create(img)
		} else {
			log.Println("Updating The Image Record")
			DB.Model(oldi).Updates(*img)
			img = oldi
		}
	} else {
		log.Printf("Saving Image %#v\n", *img)
		DB.Save(img)
	}
	log.Printf("Saved Image %#v\n", *img)
}
