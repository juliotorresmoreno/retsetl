package dictionaryGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/api/crud"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/graphql-go/graphql"
)

//remove Remove the fields of any resource in the system
func remove() *graphql.Field {
	fields := []string{"resource", "id"}
	FieldsConfigArgument := graphql.FieldConfigArgument{}
	for _, v := range fields {
		FieldsConfigArgument[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "resourceDeleteSuccess",
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
			row["id"] = helper.GetGraphParam(params, "id", "-1")
			row.SetTableName(resource.(string))
			_, err = crud.Delete(row)
			return map[string]string{"message": "OK"}, err
		},
	}
}
