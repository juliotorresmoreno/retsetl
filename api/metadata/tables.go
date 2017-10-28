package metadata

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/go-xorm/xorm"
)

//GetTables Returns the details of the tables from its id
func GetTables(conn *xorm.Engine, rows *[]models.Table, class ...string) (bool, error) {
	session := conn.NewSession()
	for _, v := range class {
		session = session.Where("class_id = ?", v)
	}
	if err := session.Find(rows); err != nil {
		return false, err
	}
	return true, nil
}

//GetTablesByID Returns the details of the tables from its id
func GetTablesByID(conn *xorm.Engine, rows *[]models.Table, ID string) (bool, error) {
	session := conn.NewSession()
	session = session.Where("id = ?", ID)
	if err := session.Find(rows); err != nil {
		return false, err
	}
	return true, nil
}

//CountClassMapFields Count the number of mapped fields in a dictionary
func CountClassMapFields(hostConn, mlsConn *xorm.Engine, resource, class string) (int, int, error) {
	mapped, err := mlsConn.
		Where("class_name = ?", class).
		Where("resource_name = ?", resource).
		Where("mapped_status = 'ALL'").
		Count(models.Table{})
	if err != nil {
		return 0, 0, err
	}
	unmapped, err := mlsConn.
		Where("class_name = ?", class).
		Where("resource_name = ?", resource).
		Where("mapped_status != 'ALL'").
		Count(models.Table{})
	if err != nil {
		return 0, 0, fmt.Errorf("Error: %v", err)
	}
	return int(mapped), int(unmapped), err
}

func toString(val interface{}) string {
	switch val.(type) {
	case []uint8:
		return string(val.([]uint8))
	default:
		return ""
	}
}

func exists(data1, data2 []string) bool {
	data1len := len(data1)
	data2len := len(data2)
	for i1 := 0; i1 < data1len; i1++ {
		for i2 := 0; i2 < data2len; i2++ {
			if data1[i1] != "" {
				if data1[i1] == data2[i2] {
					return true
				}
			}
		}
	}
	return false
}
