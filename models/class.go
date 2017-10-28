package models

import (
	"time"
)

//Class Represents the class table in the metadata of a server rets
type Class struct {
	ID            string    `xorm:"id varchar(40) not null pk"           json:"id"`
	StoreAs       string    `xorm:"store_as varchar(40) not null"        json:"store_as"`
	ResourceID    string    `xorm:"resource_id varchar(255) not null"    json:"resource_id"`
	ResourceName  string    `xorm:"resource_name varchar(255) not null"  json:"resource_name"`
	Mls           string    `xorm:"mls varchar(255) not null"            json:"mls"`
	ClassName     string    `xorm:"class_name varchar(255) not null"     json:"class_name"`
	StandardName  string    `xorm:"standard_name varchar(255) not null"  json:"standard_name"`
	VisibleName   string    `xorm:"visible_name varchar(255) not null"   json:"visible_name"`
	Description   string    `xorm:"description varchar(255) not null"    json:"description"`
	TableVersion  string    `xorm:"table_version varchar(255) not null"  json:"table_version"`
	TableDate     string    `xorm:"table_date varchar(255) not null"     json:"table_date"`
	UpdateVersion string    `xorm:"update_version varchar(255) not null" json:"update_version"`
	UpdateDate    string    `xorm:"update_date varchar(255) not null"    json:"update_date"`
	CreateAt      time.Time `xorm:"created"`
	UpdateAt      time.Time `xorm:"updated"`
}

//TableName Name of the table where the data will be stored
func (el Class) TableName() string {
	return "metadata_classnames"
}
