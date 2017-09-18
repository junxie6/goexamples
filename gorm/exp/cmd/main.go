package main

import (
	"fmt"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

import (
	"github.com/junhsieh/goexamples/gorm/hasmany"
)

var (
	db *gorm.DB
)

func main() {
	//
	var err error

	db, err = gorm.Open("mysql", "exp:exp@/exp?charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	//
	db.AutoMigrate(&model.Person{})
	db.AutoMigrate(&model.Address{})

	db.Model(&model.Address{}).AddForeignKey("id_user", "people(id)", "RESTRICT", "RESTRICT")

	// === Create address and person
	//address1 := model.Address{
	//	Street1:  "1234 abc road",
	//	City:     "Burnaby",
	//	Province: "BC",
	//	Country:  "Canada",
	//}
	//address2 := model.Address{
	//	Street1:  "456 abc road",
	//	City:     "Surrey",
	//	Province: "BC",
	//	Country:  "Canada",
	//}

	person1 := model.Person{
		FirstName: "Jun",
		LastName:  "Hsieh",
		//AddressArr: []model.Address{address1, address2},
	}

	// Create
	//db.Create(&person1)

	// Insert ot udpate
	db.Where(model.Person{FirstName: "Jun"}).Assign(model.Person{LastName: "Hsieh 2"}).FirstOrCreate(&person1)

	fmt.Printf("Person: %#v\n", person1)

	// === Get person
	//person1 := model.Person{}

	//db.First(&person1, 1)

	//fmt.Printf("Person: %#v\n", person1)
}
