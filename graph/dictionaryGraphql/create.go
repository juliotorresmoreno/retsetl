package dictionaryGraphql

import (
	"fmt"

	"github.com/go-xorm/xorm"

	"bitbucket.org/mlsdatatools/retsetl/api/crud"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"

	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/graphql-go/graphql"
)

//create Create the fields of any resource in the system
func create() *graphql.Field {
	fields := []string{
		"resource",
		"StandardName", "Definition", "Groups", "SimpleDataType",
		"SugMaxLength", "Synonym", "ElementStatus", "BEDES",
		"CertificationLevel", "RecordID", "LookupStatus", "Lookup",
		"SugMaxPrecision", "RepeatingElement", "PropertyTypes",
		"Payloads", "StatusChangeDate", "RevisedDate",
		"AddedInVersion", "Wiki", "FieldID", "FieldStatus",
	}
	FieldsConfigArgument := graphql.FieldConfigArgument{}
	for _, v := range fields {
		FieldsConfigArgument[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "resourceCreateSuccess",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args:        FieldsConfigArgument,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			var resource = helper.GetGraphParam(params, "resource", "").(string)
			var hosConn *xorm.Engine
			if resource == "" {
				return nil, fmt.Errorf("resource not found")
			}
			row := orm.Model{}
			for k, v := range params.Args {
				if k != "resource" {
					row[k] = v
				}
			}
			row.SetTableName(resource)
			if hosConn, err = db.Open(); err == nil {
				search := models.DictionaryResource{
					Table: resource,
				}
				hosConn.Sync2(search)
				cnt, _ := hosConn.Where("StandardName = ?", row["StandardName"]).Count(search)
				if cnt == 0 {
					_, err = crud.Create(row)
					return map[string]string{"message": "OK"}, err
				}
				err = fmt.Errorf("It is not possible to add the field because it already exists")
				return map[string]string{"message": "OK"}, err
			}
			return map[string]string{"message": "OK"}, err
		},
	}
}
