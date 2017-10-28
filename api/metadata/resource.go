package metadata

import (
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/go-xorm/xorm"
)

//GetResources  Returns the details of the resource from condition
func GetResources(conn *xorm.Engine, row *[]models.Resource, mls ...string) (bool, error) {
	session := conn.NewSession()
	for _, v := range mls {
		session = session.Where("mls = ?", v)
	}
	err := session.Find(row)
	return true, err
}

//GetResourcesByID Returns the details of the resource from its id
func GetResourcesByID(conn *xorm.Engine, row *models.Resource, id string) (bool, error) {
	session := conn.NewSession().Where("id = ?", id)
	return session.Get(row)
}
