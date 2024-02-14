package models

import (
	"errors"
	"mygram-api/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username     string        `gorm:"not null" json:"username" form:"username" valid:"required~Username is required"`
	Email        string        `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Email is required,email~Invalid email format"`
	Password     string        `gorm:"not null" json:"password" form:"password" valid:"required~Password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Age          string        `gorm:"not null" json:"age" form:"age" valid:"required~Age is required,`
	Photos       []Photo       `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos"`
	Comments     []Comment     `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	SocialMedias []SocialMedia `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json;"socialmedias"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	if uint(u.Age) <= 8 {
		err = errors.New("age min 8 or older")
	}

	u.Password = helpers.HashPass(u.Password)
	err = nil
	return
}
