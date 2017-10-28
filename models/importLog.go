package models

import (
	"time"

	"bitbucket.org/mlsdatatools/retsetl/db"
)

//ImportLog d
type ImportLog struct {
	ID         string    `xorm:"id pk"       json:"id"`
	Mls        string    `xorm:"mls"         json:"mls"`
	Resource   string    `xorm:"resource"    json:"resource"`
	ResourceID string    `xorm:"resource_id" json:"resource_id"`
	Class      string    `xorm:"class"       json:"class"`
	ClassID    string    `xorm:"class_id"    json:"class_id"`
	Date       time.Time `xorm:"date"        json:"date"`
	Imported   string    `xorm:"imported"    json:"imported"`
	Status     string    `xorm:"status"      json:"status"`
}

//ImportLogDetails d
type ImportLogDetails struct {
	ID          string `xorm:"id         pk" json:"id"`
	ImportLogID string `xorm:"import_log_id" json:"import_log_id"`
	Error       string `xorm:"error"         json:"error"`
	Data        string `xorm:"data"          json:"data"`
}

//TableName d
func (el ImportLog) TableName() string {
	return "import_log"
}

//TableName d
func (el ImportLogDetails) TableName() string {
	return "import_log_details"
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		conn.Sync2(ImportLog{})
		conn.Sync2(ImportLogDetails{})
	}
}
