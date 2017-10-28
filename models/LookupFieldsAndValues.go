package models

import "bitbucket.org/mlsdatatools/retsetl/db"
import "bitbucket.org/mlsdatatools/retsetl/orm"

//LookupFieldsAndValues Stores the names of available resources that can be mapped in the dictionary.
type LookupFieldsAndValues struct {
	LookupValue      string `xorm:"LookupValue"`
	LookupField      string `xorm:"LookupField"`
	Definition       string `xorm:"Definition"`
	Synonym          string `xorm:"Synonym"`
	BEDES            string `xorm:"BEDES"`
	References       string `xorm:"References"`
	LookupStatus     string `xorm:"LookupStatus"`
	LookupID         string `xorm:"LookupID"`
	LookupFieldID    string `xorm:"LookupFieldID"`
	StatusChangeDate string `xorm:"StatusChangeDate"`
	RevisedDate      string `xorm:"RevisedDate"`
	AddedInVersion   string `xorm:"AddedInVersion"`
	WikiPage         string `xorm:"WikiPage"`
	FieldID          string `xorm:"FieldID"`
	FieldStatus      string `xorm:"FieldStatus"`
}

//TableName Name of the table where the data will be stored
func (el LookupFieldsAndValues) TableName() string {
	return "LookupFieldsAndValues"
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		session := orm.NewSession(conn)
		row := orm.Model{}
		row.SetTableName(LookupFieldsAndValues{}.TableName())
		session.AddColumn(row, "FieldID", "varchar(10)")
		session.AddColumn(row, "FieldStatus", "varchar(10)")
	}
}
