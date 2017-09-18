package model

import (
	"time"
)

import (
//"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Person struct {
	Model
	FirstName  string    `gorm:"not null"`
	LastName   string    `gorm:"not null"`
	AddressArr []Address `gorm:"ForeignKey:IDUser"`
}

type Address struct {
	Model
	IDUser   uint   `gorm:"not null"`
	Street1  string `gorm:"not null"`
	Street2  string `gorm:"not null"`
	City     string `gorm:"not null"`
	Province string `gorm:"not null"`
	Country  string `gorm:"not null"`
}
