package model

type Permission struct {
	IDPermission uint `db:"IDPermission"`
	IDProject    uint `db:"IDProject"`
	IDAction     uint `db:"IDAction"`
}
