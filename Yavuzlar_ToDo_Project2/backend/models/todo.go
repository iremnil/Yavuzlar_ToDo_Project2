package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Todo struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Status string `gorm:"type:varchar(50);default:'Todo'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	//Foreign Key hangi kullanıcıya ait
	UserID uuid.UUID `json:"user_id"`
	
	//İlişki,GORMun bu Todonun birusera ait olduğunu anlaması için
	// duzetlme json:"-" ekleyerek o boş kullanıcı şablonunun API'ye sızmasının onunu kestık
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"-"` 
}
func (t *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}