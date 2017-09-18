package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Profile  Profile `gorm:"ForeignKey:UserRefer"`
	//CreditCard CreditCard
}

type CreditCard struct {
	gorm.Model
	UserID uint
	Number string
}

type Profile struct {
	gorm.Model
	Name      string
	Age       uint
	Title     string
	UserRefer uint
}
