package standards

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// Thread defines the the standard of a thread used
// todo : https://en.wikipedia.org/wiki/Screw_thread
type Thread struct {
	Standard      `gorm:"embedded" formType:"-"`
	ThreadPerInch int16   `formType:"int" formUnit:"tpi"`
	Diameter      float32 `formType:"int" formUnit:"mm"`
	Orientation   string  `formType:"string" formUnit:"orientation"`
}

func NewThread() *Thread {
	t := new(Thread)
	t.Type = "Thread"
	return t
}

// Get Thread return the requests Thread Standards ID
func (t *Thread) Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{}) {
	return t.Standard.Get(db, values, id, t, adj)
}

// Delete Thread delete the requested Thread standard ID
func (t *Thread) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return t.Standard.Delete(db, values, id, t)
}

// Post Thread delete the requested Thread standard ID
func (t *Thread) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return t.Standard.Post(db, values, request, id, adj, t)
}

// Put Thread delete the requested Thread standard ID
func (t *Thread) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return t.Standard.Put(db, values, body, t)
}

// Save Thread will register the BB into the database
func (t *Thread) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if t.GetID() == 0 {
		oldt := new(Thread)
		if errors.Is(db.Where("name = ? AND code = ?", t.GetName(), t.GetCode()).First(&oldt).Error, gorm.ErrRecordNotFound) {
			// We update our just created object in order to add it's associations ...
			err = db.Save(t).Error
			if err != nil {
				return
			}
		} else {
			err = db.Model(oldt).Updates(t).Error

			if err != nil {
				return
			}

			db.Model(oldt).First(t, oldt.ID)
		}
	} else {
		err = db.Save(t).Error
		if err != nil {
			return
		}
	}
	return
}
