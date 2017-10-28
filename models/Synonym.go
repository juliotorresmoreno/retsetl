package models

import "bitbucket.org/mlsdatatools/retsetl/db"

//Synonym Represents the mapped fields of the data dictionary
type Synonym struct {
	ID           string `xorm:"id"`
	Synonymous   string `xorm:"synonymous"`
	StandardName string `xorm:"standard_name"`
	SystemName   string `xorm:"system_name"`
	ClassName    string `xorm:"class_name"`
	ResourceName string `xorm:"resource_name"`
	FieldID      string `xorm:"field_id"`
	ClassID      string `xorm:"class_id"`
	ResourceID   string `xorm:"resource_id"`
}

//TableName Name of the table where the data will be stored
func (el Synonym) TableName() string {
	return "synonym"
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		conn.Sync2(Synonym{})
	}
}
