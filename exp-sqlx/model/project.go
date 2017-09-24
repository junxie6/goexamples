package model

import (
	"database/sql"
	//"fmt"
)

type Project struct {
	IDProject uint   `db:"IDProject"`
	Name      string `db:"Name"`
}

func ListProject() ([]Project, error) {
	var err error
	projectArr := make([]Project, 0)

	sq := "SELECT "
	sq += "  IDProject "
	sq += ", Name "
	sq += "FROM project "

	err = db.Select(&projectArr, sq)

	if err == sql.ErrNoRows {
		return projectArr, nil
	} else if err != nil {
		return nil, err
	}

	return projectArr, nil
}
