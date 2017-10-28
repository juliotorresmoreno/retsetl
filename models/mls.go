package models

import (
	"time"

	"bitbucket.org/mlsdatatools/retsetl/db"
)

//Mls Represents a server mls of some client and its database
type Mls struct {
	ID                string `xorm:"id varchar(40) not null pk"                  json:"id"                 bson:"_id"`
	Active            int    `xorm:"active bool not null"                        json:"active"             bson:"active"`
	Name              string `xorm:"name varchar(100) not null"                  json:"name"               bson:"name"               valid:"alphaspacesnum"`
	SystemID          string `xorm:"system_id varchar(100) not null"             json:"system_id"          bson:"system_id"          valid:"alphaspacesnum"`
	URL               string `xorm:"url varchar(200) not null"                   json:"url"                bson:"url"                valid:"required,url"`
	Username          string `xorm:"username varchar(100) not null"              json:"username"           bson:"username"           valid:"required"`
	Password          string `xorm:"password varchar(100) not null"              json:"password"           bson:"password"           valid:"required"`
	VersionRets       string `xorm:"version_rets varchar(100) not null"          json:"version"            bson:"version"            valid:"required,rets"`
	UseragentName     string `xorm:"useragent_name varchar(100) not null"        json:"useragent_name"     bson:"useragent_name"     valid:"utfletternum"`
	UseragentPassword string `xorm:"useragent_password varchar(100) not null"    json:"useragent_password" bson:"useragent_password" valid:""`

	ServerBd   string `xorm:"server_bd varchar(100) not null"   json:"server_bd"   bson:"server_bd"   valid:"server"`
	UsernameBd string `xorm:"username_bd varchar(100) not null" json:"username_bd" bson:"username_bd" valid:"utfletternum"`
	PasswordBd string `xorm:"password_bd varchar(100) not null" json:"password_bd" bson:"password_bd" valid:""`
	NameBd     string `xorm:"name_bd varchar(100) not null"     json:"name_bd"     bson:"name_bd"     valid:"utfletternum"`

	CreateAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"delete_at"`
}

//TableName Name of the table where the data will be stored
func (el Mls) TableName() string {
	return "mls"
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		conn.Sync2(Mls{})
	}
}
