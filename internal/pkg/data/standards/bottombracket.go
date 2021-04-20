package standards

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// BottomBracket will define the bottom bracket standard
type BottomBracket struct {
	Standard `gorm:"embedded" formType:"-"`
	// Thread definition (if needed)
	ThreadID int    `json:"-" fromType:"-"`
	Thread   Thread `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" formType:"nested"`
	// IsThreaded : true if  it's a threaded bottom bracket
	IsThreaded bool `json:"isThreaded" formType:"bool"`
	// IsPressFit : true if it's a pressfit bottom bracket
	// (can't be true with isThreaded)
	IsPressFit bool `json:"isPressFit" formType:"bool"`
	// the inside width of the shell in mm
	ShellWidth float32 `json:"shellWidth" formType:"int" formUnit:"mm"`
	// External diameter in mm
	ExternalWidth float32 `json:"externalWidth" formType:"int" formUnit:"mm"`
}

func NewBottomBracket() *BottomBracket {
	bb := new(BottomBracket)
	bb.Type = "BottomBracket"
	return bb
}

// Get BottomBracket return the requests BottomBracket Standards ID
func (bb *BottomBracket) Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{}) {
	return bb.Standard.Get(db, values, id, bb, adj)
}

// Delete BottomBracket delete the requested BottomBracket standard ID
func (bb *BottomBracket) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return bb.Standard.Delete(db, values, id, bb)
}

// Post BottomBracket delete the requested BottomBracket standard ID
func (bb *BottomBracket) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return bb.Standard.Post(db, values, request, id, adj, bb)
}

// Put BottomBracket delete the requested BottomBracket standard ID
func (bb *BottomBracket) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return bb.Standard.Put(db, values, body, bb)
}

// Save BottomBracket will register the BB into the database
func (bb *BottomBracket) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if bb.GetID() == 0 {
		oldbb := new(BottomBracket)
		if errors.Is(db.Where("name = ? AND code = ?", bb.GetName(), bb.GetCode()).First(&oldbb).Error, gorm.ErrRecordNotFound) {
			// We update our just created object in order to add it's associations ...
			err = db.Save(bb).Error
			if err != nil {
				return
			}
		} else {
			err = db.Model(oldbb).Updates(bb).Error

			if err != nil {
				return
			}

			db.Model(oldbb).First(bb, oldbb.ID)
		}
	} else {
		err = db.Save(bb).Error
		if err != nil {
			return
		}
	}
	return
}
