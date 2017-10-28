package crud

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/orm"
)

//Delete Method capable of exercising the basic functions of crud in any table
//@param row: Contains the information of the row to be treated
func Delete(row orm.Model) (bool, error) {
	conn, err := db.Open()
	if err != nil {
		return false, err
	}
	session := orm.NewSession(conn)
	err = session.Delete(row, fmt.Sprintf("id = '%v'", row["id"]))
	return err == nil, err
}
