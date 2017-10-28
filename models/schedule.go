package models

import "bitbucket.org/mlsdatatools/retsetl/db"

//Schedule Represents an import schedule
type Schedule struct {
	ID        string `xorm:"id pk"                  json:"id"`
	Mls       string `xorm:"mls       varchar(100)" json:"mls"`
	Sunday    int    `xorm:"sunday    int"          json:"sunday"`
	Monday    int    `xorm:"monday    int"          json:"monday"`
	Tuesday   int    `xorm:"tuesday   int"          json:"tuesday"`
	Wednesday int    `xorm:"wednesday int"          json:"wednesday"`
	Thursday  int    `xorm:"thursday  int"          json:"thursday"`
	Friday    int    `xorm:"friday    int"          json:"friday"`
	Saturday  int    `xorm:"saturday  int"          json:"saturday"`
	Frequency int    `xorm:"frequency int"          json:"frequency"`
	Hour      string `xorm:"hour      text"         json:"hour"`
	Classes   string `xorm:"classes   text"         json:"classes"`
	Status    string `xorm:"status    text"         json:"status"`
}

//TableName Name of the table where the data will be stored
func (el Schedule) TableName() string {
	return "schedule"
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		conn.Sync2(Schedule{})
	}
}
