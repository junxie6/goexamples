package main

// Reference:
// https://stackoverflow.com/questions/19098797/fastest-way-to-flatten-un-flatten-nested-json-objects
// https://stackoverflow.com/questions/2276463/how-can-i-get-form-data-with-javascript-jquery
// https://stackoverflow.com/questions/22985676/convert-string-with-dot-notation-to-json
// https://stackoverflow.com/questions/6393943/convert-javascript-string-in-dot-notation-into-an-object-reference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"time"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/junxie6/go-form-it"
	"github.com/junxie6/go-form-it/fields"
)

// Size constants
const (
	KB = 1 << 10
	MB = 1 << 20
)

var (
	// IOLimitReaderSize ...
	IOLimitReaderSize int64 = 2 * MB
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
	ID        uint `gorm:"primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`
}

type PGInfo struct {
	ErrorArr []string `gorm:"-";`
}

type Language struct {
}

type Profile struct {
	PGModel
	PGInfo   PGInfo
	Name     string
	Date     string
	Location string
}

type User struct {
	PGModel
	PGInfo          PGInfo
	Username        string       `gorm:"column:Username;type:varchar(32);not null;"`
	Age             uint         `gorm:"column:Age;type:varchar(32);not null;"`
	MarriageStatus  uint         `gorm:"column:MarriageStatus;not null;"`
	PrimaryLanguage uint         `gorm:"column:PrimaryLanguage;not null;"`
	CreditCardArr   []CreditCard `gorm:"ForeignKey:UserID;"`
	//Bag           Bag          `gorm:"ForeignKey:UserID;"`
	//Profile       Profile      `gorm:"ForeignKey:ProfileRefer"` // use ProfileRefer as foreign key
	//ProfileRefer  uint
}

type CreditCard struct {
	//PGModel
	PGInfo     PGInfo
	UserID     uint   `gorm:"primary_key;column:UserID;not null;" sql:"type:int(10) UNSIGNED NOT NULL DEFAULT 0"`
	Weight     uint   `gorm:"primary_key;column:Weight;not null;" sql:"type:int(10) UNSIGNED NOT NULL DEFAULT 0"`
	Number     string `gorm:"column:Number;type:varchar(32);not null;"`
	ExpireDate string `gorm:"column:ExpireDate;type:varchar(32);not null;"`
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

	//http.Handle("/static/", http.FileServer(http.Dir(".")))
	//http.HandleFunc("/", srvHome)
	//http.HandleFunc("/save", srvForm)
	//http.ListenAndServe(":8444", nil)

	Test()
}

func DropTables() {
	Conn.DropTable(&User{})
	Conn.DropTable(&Profile{})
	Conn.DropTable(&CreditCard{})
	Conn.DropTable(&Bag{})
	Conn.DropTable(&BagItem{})
}

func AutoMigrateTables() {
	Conn.AutoMigrate(&User{})
	Conn.AutoMigrate(&CreditCard{})
	//Conn.AutoMigrate(&Profile{})
	//Conn.AutoMigrate(&Bag{})
	//Conn.AutoMigrate(&BagItem{})

	//Conn.Model(&User{}).AddForeignKey("profile_refer", "profile(id)", "RESTRICT", "RESTRICT")
}

func ObjectToJSON(u1 interface{}, IsFormat bool) {
	var byteArr []byte
	var err error

	if IsFormat == true {
		byteArr, err = json.MarshalIndent(&u1, "", "    ")
	} else {
		byteArr, err = json.Marshal(&u1)
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("u1: %s\n", string(byteArr))
}

func JSONToObject(v interface{}, JSONStr string) {
	var err error

	err = json.Unmarshal([]byte(JSONStr), v)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}

func Test() {
	//var err error

	// Create profile
	//p1 := Profile{
	//	Name:     "Jun2 profile",
	//	Date:     "2018-01-25",
	//	Location: "Vancouver",
	//}
	//Conn.Save(&p1)

	// Create user
	//var u1 User

	//u1 = User{
	//	PGModel: PGModel{
	//		ID: 1,
	//	},
	//	Username: "Jun 2",
	//	Age:      18,
	//}
	//u1.CreditCardArr = []CreditCard{
	//	CreditCard{
	//		Weight:     0,
	//		Number:     "1233",
	//		ExpireDate: "2018-01-28",
	//	},
	//	CreditCard{
	//		Weight:     1,
	//		Number:     "4566",
	//		ExpireDate: "2018-01-28",
	//	},
	//}
	//Conn.Save(&u1)

	//
	//u1 = User{}
	//u1.ID = 1
	//Conn.Preload("CreditCardArr").First(&u1)

	//ObjectToJSON(u1, true)
	//ObjectToJSON(u1, false)

	var JSONStr = `
{
    "ID": 1,
    "CreatedAt": "2018-01-29T05:40:07Z",
    "UpdatedAt": "2018-01-29T05:40:07Z",
    "PGInfo": {
        "ErrorArr": null
    },
    "Username": "Jun 2",
    "Age": 18,
    "MarriageStatus": 0,
    "PrimaryLanguage": 0,
    "CreditCardArr": [
        {
            "PGInfo": {
                "ErrorArr": null
            },
            "UserID": 1,
            "Weight": 0,
            "Number": "123",
            "ExpireDate": "2018-01-28"
        },
        {
            "PGInfo": {
                "ErrorArr": null
            },
            "UserID": 1,
            "Weight": 1,
            "Number": "456",
            "ExpireDate": "2018-01-28"
        }
    ]
}
	`
	u2 := User{}
	JSONToObject(&u2, JSONStr)

	//Conn.Save(&u2)
	ObjectToJSON(u2, true)

	byteArr := GetForm(&u2)
	fmt.Printf("%s\n", string(byteArr))
}

func GetForm(v interface{}) []byte {
	opts := []fields.InputChoice{
		fields.InputChoice{Id: "User.PrimaryLanguage.111", Val: "1", Text: "PHP"},
		fields.InputChoice{Id: "User.PrimaryLanguage.222", Val: "2", Text: "Go"},
	}

	form := forms.EmptyForm(forms.POST, "/exp.html").Elements(
		fields.RadioField("User.MarriageStatus", []fields.InputChoice{
			fields.InputChoice{Id: "User.MarriageStatus.111", Val: "1", Text: "Single"},
			fields.InputChoice{Id: "User.MarriageStatus.222", Val: "2", Text: "Married"},
		}).SetLabel("Marriage Status"),
		fields.TextField("User.ID").SetLabel("ID"),
		fields.TextField("User.Username").SetLabel("Username"),
		fields.TextField("User.Age").SetLabel("Age"),
		fields.SelectField("User.PrimaryLanguage", map[string][]fields.InputChoice{
			"": opts,
		}).SetLabel("Primary Language"),
		forms.FieldSet("User.CreditCarddArr",
			fields.TextField("User.CreditCarddArr[0].ID"),
			fields.TextField("User.CreditCarddArr[0].Number"),
			fields.TextField("User.CreditCarddArr[0].ExpireDate"),
			fields.TextField("User.CreditCarddArr[1].ID"),
			fields.TextField("User.CreditCarddArr[1].Number"),
			fields.TextField("User.CreditCarddArr[1].ExpireDate"),
		),
		fields.SubmitButton("User.Save", "Save"),
	)

	form.PopulateData(v)

	var buf bytes.Buffer

	exp2TplPath := "exp.html"
	exp2Tpl := template.New(path.Base(exp2TplPath))
	exp2Tpl.ParseFiles(exp2TplPath)
	exp2Tpl.Execute(&buf, map[string]interface{}{
		"form": form,
	})

	//fmt.Printf("%#v\n", buf.String())
	return buf.Bytes()
}

func srvForm(w http.ResponseWriter, r *http.Request) {
	u1 := User{}

	if err := json.NewDecoder(io.LimitReader(r.Body, IOLimitReaderSize)).Decode(&u1); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	Conn.Save(&u1)
	//fmt.Printf("u1: %#v\n", u1)
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
			<br>User.CreditCardArr.Number: <input type="text" name="User.CreditCardArr.0.Number" />
			<br>User.Bag.Name<input type="text" name="User.Bag.Name" />
			<br><input id="myBtn" type="submit" />
		</form>

		<script src="/static/dataobject-parser/dataobject-parser.js"></script>

<script>

(function($) {
	$(document).ready(function(){
		//$('#myBtn').on('click', function(e){
		document.querySelector('#myForm').addEventListener('submit', (e) => {
			e.preventDefault();

			const formData = new FormData(e.target);

			var d = new DataObjectParser();

			var obj = {};

			for (let pair of formData.entries()) {
				//console.log(pair[0]+ ', ' + pair[1]);

				if (pair[0] == 'User.Age') {
					pair[1] = parseInt(pair[1]);
				}
				stringToObject(pair[0], pair[1], obj);


				d.set(pair[0], pair[1]);
			}

			var obj = d.data();


console.log(JSON.stringify(obj));
			return;

			var jqxhr = $.ajax({
				method: 'POST',
				url: '/save',
				data: JSON.stringify(obj.User),
				dataType: 'json',
				contentType: "application/json; charset=utf-8",
			})
			.done(function() {
			})
			.fail(function() {
			})
			.always(function() {
				console.log('hereeee');
			});
			console.log(obj.User);
			return false;
		});
	});
})(jQuery);


function stringToObject(path, value, obj) {
	var parts = path.split(".");
	var part = '';
	var last = parts.pop();

	while(part = parts.shift()) {
		if( typeof obj[part] != "object") {
			if (part == 'CreditCardArr') {
				obj[part] = [];
			} else {
				obj[part] = {};
			}
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
