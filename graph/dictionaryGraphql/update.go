package dictionaryGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/api/crud"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/graphql-go/graphql"
)

//update Update the fields of any resource in the system
func update() *graphql.Field {
	fields := []string{
		"resource", "id",
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
			Name: "resourceUpdateSuccess",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args:        FieldsConfigArgument,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			var resource = helper.GetGraphParam(params, "resource", "")
			if resource.(string) == "" {
				return nil, fmt.Errorf("resource not found")
			}
			row := orm.Model{}
			row.SetTableName(resource.(string))
			for k, v := range params.Args {
				if k != "resource" {
					row[k] = v
				}
			}
			_, err = crud.Update(row)
			return map[string]string{"message": "OK"}, err
		},
	}
}
