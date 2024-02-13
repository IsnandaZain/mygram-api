package models

type Comment struct {
	GormModel
	Message string `gorm:"not null" json:"message" form:"message"`
	UserID  uint
	User    *User
	PhotoID uint
	Photo   *Photo
}
