package crud

import (
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/go-xorm/builder"
)

//Search Method capable of exercising the basic functions of crud in any table
//@param row: Contains the information of the row to be treated
//@param cond: Search conditions
func Search(row orm.Model, cond builder.Eq) ([]map[string]interface{}, bool, error) {
	conn, err := db.Open()
	if err != nil {
		return nil, false, err
	}
	session := orm.NewSession(conn)
	data, err := session.Find(row, cond)
	return data, err == nil, err
}
