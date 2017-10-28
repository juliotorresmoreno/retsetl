package crud

import (
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/orm"
)

//Create Method capable of exercising the basic functions of crud in any table
//@param row: Contains the information of the row to be treated
func Create(row orm.Model) (bool, error) {
	conn, err := db.Open()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	session := orm.NewSession(conn)
	err = session.Insert(row)
	return err == nil, err
}
