package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// Hub holds the information on the hub
type Hub struct {
	Standard `gorm:"embedded"`
	// Holes : the number of spokes attached to the hub
	Holes uint8
	// AxleType the type of axle required
	AxleType string
	// AxelLength the Length of the Axel (external)
	AxelLength float32
	// AxelDiameter the Diameter of the axle (external)
	AxelDiameter float32
	// HasDiscBrakes : does the hub support disc brakes
	HasDiscBrakes bool
}

// Get Hub return the requests Hub Standards ID
func (h *Hub) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return h.Standard.Get(db, values, id, h)
}

// Delete Hub delete the requested Hub standard ID
func (h *Hub) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return h.Standard.Delete(db, values, id, h)
}

// Post Hub delete the requested Hub standard ID
func (h *Hub) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return h.Standard.Post(db, values, request, id, adj, h)
}

// Put Hub delete the requested Hub standard ID
func (h *Hub) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return h.Standard.Put(db, values, body, h)
}

// Save Hub will register the crank into the database
func (h *Hub) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if h.GetID() == 0 {
		oldc := new(Hub)
		if errors.Is(db.Where("name = ? AND code = ?", h.GetName(), h.GetCode()).First(&oldc).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", h)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(h).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(oldc)
			err = db.Model(oldc).Updates(h).Error
			if err != nil {
				return
			}
			log.Println(oldc)
			// TODO fix this
			db.Model(oldc).First(h, oldc.ID)
		}
	} else {
		err = db.Save(h).Error
		if err != nil {
			return
		}
	}
	return
}
