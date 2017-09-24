package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

import (
	"github.com/junhsieh/goexamples/exp-sqlx/model"
)

func RegisterRoute() {
	http.HandleFunc("/", SrvRoot)
	http.HandleFunc("/SaveUser", SrvSaveUser)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
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
			<script src="/static/jqwidgets/jqxgrid.js"></script>
			<script src="/static/jqwidgets/jqxdata.js"></script>
			<script src="/static/jqwidgets/jqxscrollbar.js"></script>
			<script src="/static/jqwidgets/jqxmenu.js"></script>
			<script src="/static/jqwidgets/jqxgrid.pager.js"></script>
			<script src="/static/jqwidgets/jqxgrid.selection.js"></script>
			<script src="/static/jqwidgets/jqxgrid.columnsresize.js"></script>
			<script src="/static/jqwidgets/jqxgrid.sort.js"></script>

			<script src="/static/user.js?%d"></script>
			<script src="/static/project.js?%d"></script>
			<script src="/static/role.js?%d"></script>
			<script src="/static/permission.js?%d"></script>
			<script src="/static/main.js?%d"></script>
		</head>
		<body>
		</body>
	</html>
	`
	now := time.Now().Unix()

	w.Write([]byte(fmt.Sprintf(html, now, now, now, now, now)))
}

func SrvSaveUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var bodyByteArr []byte

	if r.Method != "POST" {
		w.Write([]byte(fmt.Sprintf(`{"ErrMsg":"%s"}`, "http post only")))
		return
	}

	bodyByteArr, err = ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	//
	user := model.User{}

	if err = json.Unmarshal(bodyByteArr, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = user.Load(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"ErrMsg":"%s"}`, err.Error())))
		return
	}

	if user.IDUser == 0 {
		if err = user.Save(); err != nil {
			w.Write([]byte(fmt.Sprintf(`{"ErrMsg":"%s"}`, err.Error())))
			return
		}
	}

	//
	var userArr []model.User

	if userArr, err = model.ListUser(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"ErrMsg":"%s"}`, err.Error())))
		return
	}

	bodyByteArr, err = json.Marshal(map[string]interface{}{
		"Status": true,
		"Data":   userArr,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bodyByteArr)
}
