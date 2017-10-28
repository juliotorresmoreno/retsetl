package models

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/db"
)

//DictionaryResource Represents a resource within the data dictionary database.
type DictionaryResource struct {
	ID                 string `xorm:"'id' int autoincr not null pk"`
	StandardName       string `xorm:"'StandardName' TEXT not null default 'NULL'"`
	Definition         string `xorm:"'Definition' TEXT"`
	Groups             string `xorm:"'Groups' TEXT default 'NULL'"`
	SimpleDataType     string `xorm:"'SimpleDataType' TEXT default 'NULL'"`
	SugMaxLength       string `xorm:"'SugMaxLength' INTEGER default NULL"`
	Synonym            string `xorm:"'Synonym' varchar(200) default 'NULL'"`
	ElementStatus      string `xorm:"'ElementStatus' TEXT default 'NULL'"`
	BEDES              string `xorm:"'BEDES' TEXT default 'NULL'"`
	CertificationLevel string `xorm:"'CertificationLevel' TEXT default 'NULL'"`
	RecordID           string `xorm:"'RecordID' INTEGER default NULL"`
	LookupStatus       string `xorm:"'LookupStatus' TEXT default 'NULL'"`
	Lookup             string `xorm:"'Lookup' TEXT default 'NULL'"`
	SugMaxPrecision    string `xorm:"'SugMaxPrecision' INTEGER default NULL"`
	RepeatingElement   string `xorm:"'RepeatingElement' TEXT default 'NULL'"`
	PropertyTypes      string `xorm:"'PropertyTypes' TEXT"`
	Payloads           string `xorm:"'Payloads' TEXT default 'NULL'"`
	StatusChangeDate   string `xorm:"'StatusChangeDate' TEXT default 'NULL'"`
	RevisedDate        string `xorm:"'RevisedDate' TEXT default 'NULL'"`
	AddedInVersion     string `xorm:"'AddedInVersion' TEXT default 'NULL'"`
	FieldDefinition    string `xorm:"'FieldDefinition' TEXT default 'NULL'"`
	Wiki               string `xorm:"'Wiki' TEXT"`
	WikiPage           string `xorm:"'WikiPage' TEXT default 'NULL'"`
	FieldID            string `xorm:"FieldID TEXT"`
	FieldStatus        string `xorm:"FieldStatus TEXT"`
	Table              string `xorm:"-" json:"-"`
}

//TableName Name of the table where the data will be stored
func (el DictionaryResource) TableName() string {
	return el.Table
}

//SetTableName Sets the name of the table in which the data will be stored in database
func (el DictionaryResource) SetTableName(table string) {
	el.Table = table
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		resources := make([]LookupFieldsAndValues, 0)
		conn.Where("LookupValue = ?", "ResourceName").Find(&resources)

		for _, v := range resources {
			dictionary := DictionaryResource{Table: v.LookupField}
			err = conn.Sync2(dictionary)
			if err != nil {
				fmt.Println(v.LookupField, err)
			}
		}
	}
}
