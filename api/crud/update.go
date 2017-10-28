package crud

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/orm"
)

//Update Method capable of exercising the basic functions of crud in any table
//@param row: Contains the information of the row to be treated
func Update(row orm.Model) (bool, error) {
	conn, err := db.Open()
	if err != nil {
		return false, err
	}
	session := orm.NewSession(conn)
	err = session.Update(row, fmt.Sprintf("id = '%v'", fmt.Sprintf("%v", row["id"])))
	return err == nil, err
}
