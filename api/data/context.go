package data

import (
	"fmt"
	"strings"

	"bitbucket.org/mlsdatatools/retsetl/orm"
)

func createContext(session *orm.Session, row orm.Model, columns []columns) {
	/*var err error
	tableName := row.TableName()
	switch tableName {
	case "Contact":
		err = createTableContact(session.Engine(), tableName)
	case "HistoryTransactional":
		err = createTableHistoryTransactional(session.Engine(), tableName)
	case "Media":
		err = createTableMedia(session.Engine(), tableName)
	case "Member":
		err = createTableMember(session.Engine(), tableName)
	case "Office":
		err = createTableOffice(session.Engine(), tableName)
	case "OpenHouse":
		err = createTableOpenHouse(session.Engine(), tableName)
	case "Property":
		err = createTableProperty(session.Engine(), tableName)
	default:
	}
	if err != nil && !strings.Contains(err.Error(), "exists") {
		fmt.Println(err)
		return
	}*/
	err := session.CreateTable(row, "postgres")
	if err != nil && !strings.Contains(err.Error(), "exists") {
		fmt.Println(err)
		return
	}
	session.AddColumn(row, "id", "serial primary key")
	for _, v := range columns {
		session.AddColumn(row, v.SystemName, "")
	}
	session.AddColumn(row, "Source", "")
	session.AddColumn(row, "UpdateDate", "")

}
