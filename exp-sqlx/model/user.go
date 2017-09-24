package model

import (
	"database/sql"
)

type User struct {
	IDUser   uint   `db:"IDUser"`
	Username string `db:"Username"`
}

func (u *User) Load() error {
	var err error

	sq := "SELECT "
	sq += "  IDUser "
	sq += ", Username "
	sq += "FROM user "
	sq += "WHERE "
	sq += "Username = ? "

	err = db.Get(u, sq, u.Username)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (u *User) Save() error {
	var err error

	_, err = db.NamedExec(`INSERT INTO user (Username) VALUES (:Username)`,
		map[string]interface{}{
			"Username": u.Username,
		})

	return err
}
