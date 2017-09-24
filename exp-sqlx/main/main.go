package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	var err error

	http.HandleFunc("/", SrvRoot)
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}
}

func SrvUser(w http.ResponseWriter, r *http.Request) {
	//	var err error
	//	var bodyByteArr []byte
	//
	//	bodyByteArr, err = ioutil.ReadAll(r.Body)
	//
	//	if err != nil {
	//		http.Error(w, "can't read body", http.StatusBadRequest)
	//		return
	//	}
	//
	//	user := model.User{}
	//
	//	err = json.Unmarshal(bodyByteArr, &user)
	//
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//		return
	//	}
	//
	//	if user.ID > 0 {
	//		//db.First(&user, user.ID)
	//
	//		db.Save(&user)
	//	} else {
	//		db.Create(&user)
	//	}
	//
	//	bodyByteArr, err = json.Marshal(&user)
	//
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//		return
	//	}
	//
	//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//	w.Write(bodyByteArr)
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
