package models

type User struct {
	GormModel
	Username     string        `gorm:"not null" json:"username" form:"username"`
	Email        string        `gorm:"not null;uniqueIndex" json:"email" form:"email"`
	Password     string        `gorm:"not null" json:"password" form:"password"`
	Age          string        `gorm:"not null" json:"age" form:"age"`
	Photos       []Photo       `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos"`
	Comments     []Comment     `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	SocialMedias []SocialMedia `gorm:"constaint:OnUpdate:CASCADE,OnDelete:SET NULL;" json;"socialmedias"`
}
