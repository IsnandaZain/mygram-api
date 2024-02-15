package models

import (
	"gorm.io/gorm"
)

type Photo struct {
	GormModel
	Title    string    `gorm:"not null" json:"title" form:"title"`
	Caption  string    `json:"caption" form:"caption"`
	PhotoUrl string    `gorm:"not null" json:"photo_url" form:"photo_url"`
	Comments []Comment `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	UserID   uint
	User     *User
}

func (p *Photo) GrabUserData(tx *gorm.DB) {
	tx.First(&p.User, "id=?", p.UserID)
}
