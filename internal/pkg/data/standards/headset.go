package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// Headset, holds the specs for a Front Derailleur
type Headset struct {
	Standard `gorm:"embedded" formType:"-"`
	// SHISAnnotation the actual S.H.I.S Annotation like ZS44/28.6"	SHISAnnotation string
	SHISAnnotation string `formType:"string"`
	// TODO : be able to list the possible values
	// FitType , is headset integrated or PressFit ?
	FitType string `formType:"String"`
}

func NewHeadset() *Headset {
	hs := new(Headset)
	hs.Type = "Headset"
	return hs
}

// Get Headset return the requests Headset Standards ID
func (rd *Headset) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return rd.Standard.Get(db, values, id, rd)
}

// Delete Headset delete the requested Headset standard ID
func (rd *Headset) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return rd.Standard.Delete(db, values, id, rd)
}

// Post Headset delete the requested Headset standard ID
func (rd *Headset) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return rd.Standard.Post(db, values, request, id, adj, rd)
}

// Put Headset delete the requested Headset standard ID
func (rd *Headset) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return rd.Standard.Put(db, values, body, rd)
}

// Save Headset will register the crank into the database
func (rd *Headset) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if rd.GetID() == 0 {
		oldc := new(Headset)
		if errors.Is(db.Where("name = ? AND code = ?", rd.GetName(), rd.GetCode()).First(&oldc).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", rd)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(rd).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(oldc)
			err = db.Model(oldc).Updates(rd).Error
			if err != nil {
				return
			}
			log.Println(oldc)
			// TODO fix this
			db.Model(oldc).First(rd, oldc.ID)
		}
	} else {
		err = db.Save(rd).Error
		if err != nil {
			return
		}
	}
	return
}
