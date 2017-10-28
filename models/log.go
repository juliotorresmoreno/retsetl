package models

import (
	"time"

	"bitbucket.org/mlsdatatools/retsetl/db"
)

//Log Name of the table where the import history of the data is stored
type Log struct {
	ID           string    `xorm:"id"            json:"id"`
	Date         time.Time `xorm:"date datetime" json:"date"`
	ResourceName string    `xorm:"resource_name" json:"resource_name"`
	Mls          string    `xorm:"mls"           json:"mls"`
	ClassID      string    `xorm:"class_id"      json:"class_id"`
	Action       string    `xorm:"action" valid:"in(transform)" json:""`
}

//LogDetails Contains the details of the imported rows
type LogDetails struct {
	ID           string          `xorm:"id"            json:"id"`
	Mls          string          `xorm:"mls"           json:"mls"`
	LogID        string          `xorm:"log_id"        json:"log_id"`
	Synonymous   bool            `xorm:"synonymous"    json:"synonymous"`
	MaxLength    bool            `xorm:"max_length"    json:"max_length"`
	Type         bool            `xorm:"type"          json:"type"`
	StandardName string          `xorm:"standard_name" json:"standard_name"`
	SystemName   string          `xorm:"system_name"   json:"system_name"`
	ResourceName string          `xorm:"resource_name" json:"resource_name"`
	FieldID      string          `xorm:"field_id"      json:"field_id"`
	ClassID      string          `xorm:"class_id"      json:"class_id"`
	ResourceID   string          `xorm:"resource_id"   json:"resource_id"`
	Status       string          `xorm:"status"        json:"status"`
	Compare      map[string]bool `xorm:"compare"       json:"compare"`
}

//TableName Name of the table where the data will be stored
func (el Log) TableName() string {
	return "log"
}

//TableName Name of the table where the data will be stored
func (el LogDetails) TableName() string {
	return "log_details"
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		conn.Sync2(Log{})
		conn.Sync2(LogDetails{})
	}
}
