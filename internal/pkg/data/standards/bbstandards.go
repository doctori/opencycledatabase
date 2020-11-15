package standards

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

// BBStandard will define the bottom bracket standard
type BBStandard struct {
	Standard `gorm:"embedded;embeddedPrefix:bb_"`
	// Thread definition (if needed)
	ThreadID int            `json:"-"`
	Thread   ThreadStandard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// IsThreaded : true if  it's a threaded bottom bracket
	IsThreaded bool
	// IsPressFit : true if it's a pressfit bottom bracket
	// (can't be true with isThreaded)
	IsPressFit bool
	// the inside width of the shell in mm
	ShellWidth float32
	// External diameter in mm
	ExternalWidth float32
}

// Post will save the BBStandard
func (BBStandard) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	body := request.Body
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	var standard BBStandard
	err := decoder.Decode(&standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)
	}
	if standard.isNul() {
		return http.StatusBadRequest, "The object is null"
	}
	log.Println(standard)
	err = standard.save(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could Not Save the Standard : \n\t %s", err.Error())
	}
	return http.StatusAccepted, standard
}

// Put updates BBStandard
func (BBStandard) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	var standard BBStandard
	err := decoder.Decode(&standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)

	}
	log.Println(standard)
	err = standard.save(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could Not Save the Standard : \n\t%s", err.Error())
	}
	return http.StatusOK, standard
}

// GetAll will returns al the standards
func GetAll(db *gorm.DB, page string, perPage string) (standards []BBStandard) {
	ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iperPage, err := strconv.Atoi(perPage)
	if err != nil {
		iperPage = defaultPerPage
	}
	//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
	db.Order("name").Offset(ipage * iperPage).Limit(iperPage).Find(&standards)
	return
}

// Get BBStandard return the requests BB Standards
func (BBStandard) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	page := values.Get("page")
	perPage := values.Get("per_page")
	log.Printf("Id is : %d\n", id)
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == 0 {
		return 200, GetAll(db, page, perPage)
	}

	var standard BBStandard
	err := db.First(&standard, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Standard not found"
	}
	return 200, standard
}

// Delete BB standard will remove the BB sTandard struct in the database
func (s *BBStandard) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	// TODO : implement Delete method
	db.Delete(s)
	return 200, ""
}

func (s *BBStandard) save(db *gorm.DB) (err error) {
	// If we have a new record we create it
	if s.ID == 0 {
		olds := new(BBStandard)
		if errors.Is(db.Where("name = ? AND code = ?", s.Name, s.Code).First(&olds).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", s)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(s).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(olds)
			err = db.Model(olds).Updates(s).Error
			if err != nil {
				return
			}
			log.Println(olds)
			*s = *olds
		}
	} else {
		err = db.Save(&s).Error
		if err != nil {
			return
		}
	}
	return
}
