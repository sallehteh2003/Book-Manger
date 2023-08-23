package Database

import (
	"gorm.io/gorm"
	"time"
)

type UserDb struct {
	gorm.Model
	Username    string `gorm:"varchar(50),unique"`
	Firstname   string `gorm:"varchar(30)"`
	Lastname    string `gorm:"varchar(30)"`
	Email       string `gorm:"varchar(50),unique"`
	Password    string `gorm:"varchar(50)"`
	PhoneNumber string `gorm:"varchar(20),unique"`
	Gender      string `gorm:"varchar(10)"`
	Books       []Book `gorm:"foreignKey:UserDbID"`
}
type Book struct {
	gorm.Model
	Name            string `gorm:"varchar(30)"`
	Category        string `gorm:"varchar(30)"`
	Volume          int
	PublishedAt     time.Time
	Summary         string `gorm:"varchar(50)"`
	TableOfContents string `gorm:"varchar(50)"`
	Publisher       string `gorm:"varchar(50)"`
	Author          Author
	UserDbID        uint
}
type Author struct {
	gorm.Model
	BookID      uint
	Firstname   string
	Lastname    string
	Birthday    time.Time
	Nationality string
}
