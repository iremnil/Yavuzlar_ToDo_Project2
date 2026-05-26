package models

import (
    "github.com/google/uuid"
    "gorm.io/gorm" 
)

type User struct {
    ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
    Username string    `gorm:"unique;not null" json:"username"`
    Password string    `json:"password"`
    
    //kullanıcıya ait tüm todolar burada duracak
    Todos    []Todo    `gorm:"foreignKey:UserID;references:ID"`
}
// kullanıcı kaydettiğim an UUID'yi kendi kendine oluştursun
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}