package models

import (
	"time"
)

//Resource Represents the resource table in the metadata of a server rets
type Resource struct {
	ID                          string    `xorm:"id varchar(40) not null pk"                         json:"id"                            bson:"_id"`
	Mls                         string    `xorm:"mls varchar(40) not null"                           json:"mls"`
	StoreAs                     string    `xorm:"store_as varchar(40) not null"                      json:"store_as"                      bson:"store_as"`
	ResourceID                  string    `xorm:"resource_id varchar(40) not null"                   json:"resource_id"                   bson:"resource_id"`
	ResourceName                string    `xorm:"resource_name varchar(40) not null"                 json:"resource_name"                 bson:"resource_name"`
	StandardName                string    `xorm:"standard_name varchar(40) not null"                 json:"standard_name"                 bson:"standard_name"`
	VisibleName                 string    `xorm:"visible_name varchar(40) not null"                  json:"visible_name"                  bson:"visible_name"`
	Description                 string    `xorm:"description varchar(40) not null"                   json:"description"                   bson:"description"`
	KeyField                    string    `xorm:"key_field varchar(40) not null"                     json:"key_field"                     bson:"key_field"`
	ClassCount                  string    `xorm:"class_count varchar(40) not null"                   json:"class_count"                   bson:"class_count"`
	ClassVersion                string    `xorm:"class_version varchar(40) not null"                 json:"class_version"                 bson:"class_version"`
	ClassDate                   string    `xorm:"class_date varchar(40) not null"                    json:"class_date"                    bson:"class_date"`
	ObjectVersion               string    `xorm:"object_version varchar(40) not null"                json:"object_version"                bson:"object_version"`
	ObjectDate                  string    `xorm:"object_date varchar(40) not null"                   json:"object_date"                   bson:"object_date"`
	SearchHelpVersion           string    `xorm:"search_help_version varchar(40) not null"           json:"search_help_version"           bson:"search_help_version"`
	SearchHelpDate              string    `xorm:"search_help_date varchar(40) not null"              json:"search_help_date"              bson:"search_help_date"`
	EditMaskVersion             string    `xorm:"edit_mask_version varchar(40) not null"             json:"edit_mask_version"             bson:"edit_mask_version"`
	EditMaskDate                string    `xorm:"edit_mask_date varchar(40) not null"                json:"edit_mask_date"                bson:"edit_mask_date"`
	LookupVersion               string    `xorm:"lookup_version varchar(40) not null"                json:"lookup_version"                bson:"lookup_version"`
	LookupDate                  string    `xorm:"lookup_date varchar(40) not null"                   json:"lookup_date"                   bson:"lookup_date"`
	UpdateHelpVersion           string    `xorm:"update_help_version varchar(40) not null"           json:"update_help_version"           bson:"update_help_version"`
	UpdateHelpDate              string    `xorm:"update_help_date varchar(40) not null"              json:"update_help_date"              bson:"update_help_date"`
	ValidationExpressionVersion string    `xorm:"validation_expression_version varchar(40) not null" json:"validation_expression_version" bson:"validation_expression_version"`
	ValidationExpressionDate    string    `xorm:"validation_expression_date varchar(40) not null"    json:"validation_expression_date"    bson:"validation_expression_date"`
	ValidationLookupVersion     string    `xorm:"validation_lookup_version varchar(40) not null"     json:"validation_lookup_version"     bson:"validation_lookup_version"`
	ValidationLookupDate        string    `xorm:"validation_lookup_date varchar(40) not null"        json:"validation_lookup_date"        bson:"validation_lookup_date"`
	ValidationExternalVersion   string    `xorm:"validation_external_version varchar(40) not null"   json:"validation_external_version"   bson:"validation_external_version"`
	ValidationExternalDate      string    `xorm:"validation_external_date varchar(40) not null"      json:"validation_external_date"      bson:"validation_external_date"`
	CreateAt                    time.Time `xorm:"created"`
	UpdateAt                    time.Time `xorm:"updated"`
}

//TableName Name of the table where the data will be stored
func (el Resource) TableName() string {
	return "metadata_resources"
}
