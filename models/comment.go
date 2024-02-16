package models

// Comment represents the model for an comment
type Comment struct {
	GormModel
	Message string `gorm:"not null" json:"message" form:"message"`
	UserID  uint
	User    *User
	PhotoID uint
	Photo   *Photo
}
