package metadata

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/go-xorm/xorm"
)

//GetClass Returns the details of the class from condition
func GetClass(conn *xorm.Engine, row *[]models.Class, cond map[string]interface{}) (bool, error) {
	session := conn.NewSession()
	defer session.Close()
	for k, v := range cond {
		session = session.Where(fmt.Sprintf("%v = ?", k), v)
	}
	err := session.Find(row)
	return true, err
}

//GetClassByID Returns the details of the class from its id resource
func GetClassByID(conn *xorm.Engine, row *models.Class, id ...string) (bool, error) {
	session := conn.NewSession()
	defer session.Close()
	for _, v := range id {
		session = session.Where("id = ?", v)
	}
	return session.Get(row)
}
