package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

type Wheel struct {
	Standard `gorm:"embedded"`
	// Diameter : the Diameter of the wheel (let's say mm)
	Diameter int16 `formType:"int" formUnit:"mm"`
}

func NewWheel() *Wheel {
	w := new(Wheel)
	w.Type = "Wheel"
	return w
}

// Get Wheel return the requests Wheel Standards ID
func (w *Wheel) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return w.Standard.Get(db, values, id, w)
}

// Delete Wheel delete the requested Wheel standard ID
func (w *Wheel) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return w.Standard.Delete(db, values, id, w)
}

// Post Wheel delete the requested Wheel standard ID
func (w *Wheel) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return w.Standard.Post(db, values, request, id, adj, w)
}

// Put Wheel delete the requested Wheel standard ID
func (w *Wheel) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return w.Standard.Put(db, values, body, w)
}

// Save Wheel will register the crank into the database
func (w *Wheel) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if w.GetID() == 0 {
		oldc := new(Wheel)
		if errors.Is(db.Where("name = ? AND code = ?", w.GetName(), w.GetCode()).First(&oldc).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", w)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(w).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(oldc)
			err = db.Model(oldc).Updates(w).Error
			if err != nil {
				return
			}
			log.Println(oldc)
			// TODO fix this
			db.Model(oldc).First(w, oldc.ID)
		}
	} else {
		err = db.Save(w).Error
		if err != nil {
			return
		}
	}
	return
}
