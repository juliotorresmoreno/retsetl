package dictionaryGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/db"

	"bitbucket.org/mlsdatatools/retsetl/models"

	"bitbucket.org/mlsdatatools/retsetl/api/crud"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/go-xorm/builder"
	"github.com/graphql-go/graphql"
)

func query() *graphql.Field {
	fields := []string{
		"id",
		"StandardName", "Definition", "Groups", "SimpleDataType",
		"SugMaxLength", "Synonym", "ElementStatus", "BEDES",
		"CertificationLevel", "RecordID", "LookupStatus", "Lookup",
		"SugMaxPrecision", "RepeatingElement", "PropertyTypes",
		"Payloads", "StatusChangeDate", "RevisedDate",
		"AddedInVersion", "Wiki", "FieldID", "FieldStatus",
	}
	Fields := graphql.Fields{}
	FieldsConfigArgument := graphql.FieldConfigArgument{
		"resource": &graphql.ArgumentConfig{Type: graphql.String},
		"mls":      &graphql.ArgumentConfig{Type: graphql.String},
		"class_id": &graphql.ArgumentConfig{Type: graphql.String},
	}
	for _, v := range fields {
		Fields[v] = &graphql.Field{Type: graphql.String}
		FieldsConfigArgument[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	return &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name:   "QuerySuccess",
			Fields: Fields,
		})),
		Description: "",
		Args:        FieldsConfigArgument,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			var data []map[string]interface{}
			var resource = helper.GetGraphParam(params, "resource", "").(string)
			var classID = helper.GetGraphParam(params, "class_id", "").(string)
			var mls = helper.GetGraphParam(params, "mls", "").(string)
			if resource == "" {
				return nil, fmt.Errorf("resource not found")
			}
			row := orm.Model{}
			row.SetTableName(resource)
			cond := builder.Eq{}
			for k, v := range params.Args {
				if k != "resource" && k != "class_id" && k != "mls" {
					cond[k] = v.(string)
				}
			}
			data, _, err = crud.Search(row, cond)
			if mls != "" {
				fields := make([]models.Table, 0)
				mlsConn, _ := db.OpenMls(mls)
				err = mlsConn.Where("mapped_status != '' and class_id = ?", classID).Find(&fields)
				if err != nil {
					return data, err
				}
				for k, row := range data {
					data[k]["FieldID"] = ""
					data[k]["FieldStatus"] = ""
					for _, field := range fields {
						if row["StandardName"] == field.Mapped {
							data[k]["FieldID"] = field.ID
							data[k]["FieldStatus"] = field.MappedStatus
						}
					}
				}
			}
			return data, err
		},
	}
}
