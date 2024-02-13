package models

type SocialMedia struct {
	GormModel
	Name   string `gorm:"not null" json:"name" form:"name"`
	UserID uint
	User   *User
}
