package main

import (
	"fmt"
	"net/http"
	"time"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	var conn *gorm.DB
	var err error

	conn, err = gorm.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	// Setting
	conn.DB().SetMaxIdleConns(10)
	conn.DB().SetMaxOpenConns(100)
	conn.SingularTable(true)
	conn.LogMode(true)

	return conn, nil
}

var (
	Conn *gorm.DB
)

type PGModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`
}

type PGInfo struct {
	ErrorArr []string `gorm:"-";`
}

type Language struct {
}

type User struct {
	PGModel
	PGInfo        PGInfo
	Name          string       `gorm:"column:Name;type:varchar(32);not null;"`
	Age           uint         `gorm:"column:Age;type:varchar(32);not null;"`
	CreditCardArr []CreditCard `gorm:"ForeignKey:UserID;"`
	Bag           Bag          `gorm:"ForeignKey:UserID;"`
}

type CreditCard struct {
	PGModel
	PGInfo PGInfo
	UserID uint   `gorm:"column:UserID;not null;"`
	Number string `gorm:"column:Number;type:varchar(32);not null;"`
}

type Bag struct {
	PGModel
	PGInfo     PGInfo
	UserID     uint      `gorm:"column:UserID;not null"`
	Name       string    `gorm:"column:Name;type:varchar(32);not null;"`
	BagItemArr []BagItem `gorm:"ForeignKey:BagID;"`
}

type BagItem struct {
	PGModel
	PGInfo PGInfo
	BagID  uint   `gorm:"column:BagID;not null;"`
	Name   string `gorm:"column:Name;type:varchar(32);not null;"`
}

func main() {
	//
	var err error

	Conn, err = Connect("exp:exp@tcp(localhost:3306)/exp?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	defer Conn.Close()

	// Drop the schemas
	//DropTables()

	// Migrate the schemas
	//AutoMigrateTables()

	http.HandleFunc("/", srvHome)
	http.HandleFunc("/save", srvForm)
	http.ListenAndServe(":8444", nil)
}

func DropTables() {
	Conn.DropTable(&User{})
	Conn.DropTable(&CreditCard{})
	Conn.DropTable(&Bag{})
	Conn.DropTable(&BagItem{})
}

func AutoMigrateTables() {
	Conn.AutoMigrate(&User{})
	Conn.AutoMigrate(&CreditCard{})
	Conn.AutoMigrate(&Bag{})
	Conn.AutoMigrate(&BagItem{})
}

func srvForm(w http.ResponseWriter, r *http.Request) {
}

func srvHome(w http.ResponseWriter, r *http.Request) {
	var html = `<!DOCTYPE html>
<html>
	<head>
		<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
	</head>
	<body>
		<form id="myForm" action="/save" method="POST">
			<br>User.Name: <input type="text" name="User.Name" />
			<br>User.Age: <input type="text" name="User.Age" />
			<br>User.CreditCardArr.Number: <input type="text" name="User.CreditCardArr.Number" />
			<br>User.Bag.Name<input type="text" name="User.Bag.Name" />
			<br><input type="submit" />
		</form>

<script>
document.querySelector('#myForm').addEventListener('submit', (e) => {
	e.preventDefault();

	const formData = new FormData(e.target);
	var obj = {};

	for (let pair of formData.entries()) {
		//console.log(pair[0]+ ', ' + pair[1]); 

		stringToObject(pair[0], pair[1], obj);
	}

	console.log(obj);
});

// Reference:
// https://stackoverflow.com/questions/22985676/convert-string-with-dot-notation-to-json
// https://stackoverflow.com/questions/2276463/how-can-i-get-form-data-with-javascript-jquery
function stringToObject(path, value, obj) {
	var parts = path.split(".");
	var part = '';
	var last = parts.pop();

	while(part = parts.shift()) {
		if( typeof obj[part] != "object") {
			obj[part] = {}
		};
		obj = obj[part]; // update "pointer"
	}
	obj[last] = value;
}
</script>
	</body>
</html>
	`
	w.Write([]byte(html))
}
