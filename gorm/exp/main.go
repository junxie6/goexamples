// go get -u github.com/jinzhu/gorm
// go get -u github.com/go-sql-driver/mysql
//
// GRANT ALL PRIVILEGES ON exp.* TO 'exp'@'localhost' IDENTIFIED BY 'exp';
// FLUSH PRIVILEGES;
//
// SHOW GLOBAL VARIABLES WHERE variable_name REGEXP 'general_log|log_output';
// #SET GLOBAL general_log = 'ON';
// #SET GLOBAL general_log = 'OFF';
// #SET GLOBAL general_log_file = '/var/lib/mysql/ubun-gui.log';
// #SET GLOBAL general_log_file = 'table';
// #SET GLOBAL log_output = 'TABLE';
//
// SELECT
// event_time, command_type, CONVERT(argument USING utf8)
// FROM mysql.general_log
// WHERE event_time >= (NOW() + INTERVAL -10 SECOND) AND command_type IN ('Execute')
// ORDER BY event_time
//
// https://facebook.github.io/react/docs/installation.html#creating-a-new-application
//
// apt-get purge nodejs npm
// curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
//
// npm install -g create-react-app
// create-react-app my-app
// cd my-app
// npm start
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

import (
	"github.com/junhsieh/goexamples/gorm/hasone"
)

var (
	db *gorm.DB
)

func main() {
	var err error

	db, err = gorm.Open("mysql", "exp:exp@/exp?charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	//
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Profile{})

	http.HandleFunc("/user", SrvUser)
	http.HandleFunc("/", SrvRoot)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err.Error())
	}
}

func SrvUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var bodyByteArr []byte

	bodyByteArr, err = ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	user := model.User{}

	err = json.Unmarshal(bodyByteArr, &user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.ID > 0 {
		//db.First(&user, user.ID)

		db.Save(&user)
	} else {
		db.Create(&user)
	}

	bodyByteArr, err = json.Marshal(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bodyByteArr)
}

func SrvRoot(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
	<html>
		<head>
			<link rel="stylesheet" href="/static/bootstrap/bootstrap.min.css" type="text/css" />
			<link rel="stylesheet" href="/static/jqwidgets/styles/jqx.base.css" type="text/css" />
			<script src="/static/jquery-2.2.4.min.js"></script>
			<script src="/static/jqwidgets/jqxcore.js"></script>
			<script src="/static/jqwidgets/jqxwindow.js"></script>
			<script src="/static/jqwidgets/jqxbuttons.js"></script>
			<script src="/static/main.js?%d"></script>
		</head>
		<body>
		</body>
	</html>
	`
	w.Write([]byte(fmt.Sprintf(html, time.Now().Unix())))
}
