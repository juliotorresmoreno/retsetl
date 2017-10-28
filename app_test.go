package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-xorm/builder"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/orm"
)

func TestOrm(t *testing.T) {
	var err error
	value1 := orm.Model{"field1": "value1"}
	value2 := orm.Model{"field1": "value1"}
	value1.SetTableName("testing")
	value2.SetTableName("testing")
	conn := db.GetConnection()
	session := orm.NewSession(conn)
	defer session.Close()
	if err = session.CreateTable(value1); err != nil {
		if !strings.Contains(err.Error(), "exists") {
			t.Error(err)
		}
	}
	if err = session.AddColumn(value1, "field1"); err != nil {
		t.Error(err)
	}
	if err = session.Insert(value1); err != nil {
		t.Error(err)
	}
	if err = session.Insert(value2); err != nil {
		t.Error(err)
	}
	value1["field1"] = "value2"
	if err = session.Update(value1, "id='"+value1["id"].(string)+"'"); err != nil {
		t.Error(err)
	}
	if err = session.Delete(value2, "id='"+value2["id"].(string)+"'"); err != nil {
		t.Error(err)
	}
	result, err := session.Find(value1, builder.Eq{"id": value1["id"]})
	if err != nil {
		t.Error(err)
	}
	if err = session.Delete(value2, "id='"+value2["id"].(string)+"'"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}
