package orm

import (
	"fmt"
	"testing"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"
)

func TestGetFields(t *testing.T) {
	conn := db.GetConnection()
	defer conn.Close()
	session := NewSession(conn)
	model := Model{}
	model.SetTableName(models.DictionaryTable{}.TableName())
	fields := GetFieldsDB(session, model)
	fmt.Println(fields)
}
