package rets

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/models"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"bitbucket.org/mlsdatatools/retsetl/parse"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2/bson"
)

func importClass(conn *xorm.Engine, mls string, meta parse.Meta, resources map[string]string) (map[string]map[string]string, []error) {
	conn.Delete(models.Class{Mls: mls})
	result := map[string]map[string]string{}
	errors := make([]error, 0)
	data := make([]models.Class, 0)
	for _, v := range meta.Rets.MetadataClass {
		resource, ok := resources[v.Resource]
		result[v.Resource] = map[string]string{"id": resource}
		if ok {
			length := len(v.Data)
			for k := 0; k < length; k++ {
				row := models.Class{
					ID:            bson.NewObjectId().Hex(),
					ResourceID:    resource,
					ResourceName:  v.Resource,
					Mls:           mls,
					ClassName:     v.Get(k, "ClassName"),
					StandardName:  v.Get(k, "StandardName"),
					VisibleName:   v.Get(k, "VisibleName"),
					Description:   v.Get(k, "Description"),
					TableVersion:  v.Get(k, "TableVersion"),
					TableDate:     v.Get(k, "TableDate"),
					UpdateVersion: v.Get(k, "UpdateVersion"),
					UpdateDate:    v.Get(k, "UpdateDate"),
				}
				data = append(data, row)
				_, err := conn.Insert(row)
				result[v.Resource][row.ClassName] = row.ID
				if err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	return result, errors
}

var classTable = models.Class{}.TableName()

func ClassStoreAs(conn *xorm.Engine, classID, StoreAs string) error {
	row := orm.Model{}
	row.SetTableName(classTable)
	row["id"] = classID
	row["store_as"] = StoreAs
	session := orm.NewSession(conn)
	err := session.Update(row, fmt.Sprintf("id = '%v'", classID))
	return err
}
