package metadataGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/api/metadata"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/graphql-go/graphql"
)

//GetData Contains the anchor points graphql on which to consult the data
var GetData = graphql.Fields{
	"metadataResources": &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "MetadataResourceItem",
			Fields: graphql.Fields{
				"id":                            &graphql.Field{Type: graphql.String},
				"resource_id":                   &graphql.Field{Type: graphql.String},
				"standard_name":                 &graphql.Field{Type: graphql.String},
				"visible_name":                  &graphql.Field{Type: graphql.String},
				"description":                   &graphql.Field{Type: graphql.String},
				"key_field":                     &graphql.Field{Type: graphql.String},
				"class_count":                   &graphql.Field{Type: graphql.String},
				"class_version":                 &graphql.Field{Type: graphql.String},
				"class_date":                    &graphql.Field{Type: graphql.String},
				"object_version":                &graphql.Field{Type: graphql.String},
				"object_date":                   &graphql.Field{Type: graphql.String},
				"search_help_version":           &graphql.Field{Type: graphql.String},
				"search_help_date":              &graphql.Field{Type: graphql.String},
				"edit_mask_version":             &graphql.Field{Type: graphql.String},
				"edit_mask_date":                &graphql.Field{Type: graphql.String},
				"lookup_version":                &graphql.Field{Type: graphql.String},
				"lookup_date":                   &graphql.Field{Type: graphql.String},
				"update_help_version":           &graphql.Field{Type: graphql.String},
				"update_help_date":              &graphql.Field{Type: graphql.String},
				"validation_expression_version": &graphql.Field{Type: graphql.String},
				"validation_expression_date":    &graphql.Field{Type: graphql.String},
				"validation_lookup_version":     &graphql.Field{Type: graphql.String},
				"validation_lookup_date":        &graphql.Field{Type: graphql.String},
				"validation_external_version":   &graphql.Field{Type: graphql.String},
				"validation_external_date":      &graphql.Field{Type: graphql.String},
			},
		})),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":                           &graphql.ArgumentConfig{Type: graphql.String},
			"id":                            &graphql.ArgumentConfig{Type: graphql.String},
			"resource_id":                   &graphql.ArgumentConfig{Type: graphql.String},
			"standard_name":                 &graphql.ArgumentConfig{Type: graphql.String},
			"visible_name":                  &graphql.ArgumentConfig{Type: graphql.String},
			"description":                   &graphql.ArgumentConfig{Type: graphql.String},
			"key_field":                     &graphql.ArgumentConfig{Type: graphql.String},
			"class_count":                   &graphql.ArgumentConfig{Type: graphql.String},
			"class_version":                 &graphql.ArgumentConfig{Type: graphql.String},
			"class_date":                    &graphql.ArgumentConfig{Type: graphql.String},
			"object_version":                &graphql.ArgumentConfig{Type: graphql.String},
			"object_date":                   &graphql.ArgumentConfig{Type: graphql.String},
			"search_help_version":           &graphql.ArgumentConfig{Type: graphql.String},
			"search_help_date":              &graphql.ArgumentConfig{Type: graphql.String},
			"edit_mask_version":             &graphql.ArgumentConfig{Type: graphql.String},
			"edit_mask_date":                &graphql.ArgumentConfig{Type: graphql.String},
			"lookup_version":                &graphql.ArgumentConfig{Type: graphql.String},
			"lookup_date":                   &graphql.ArgumentConfig{Type: graphql.String},
			"update_help_version":           &graphql.ArgumentConfig{Type: graphql.String},
			"update_help_date":              &graphql.ArgumentConfig{Type: graphql.String},
			"validation_expression_version": &graphql.ArgumentConfig{Type: graphql.String},
			"validation_expression_date":    &graphql.ArgumentConfig{Type: graphql.String},
			"validation_lookup_version":     &graphql.ArgumentConfig{Type: graphql.String},
			"validation_lookup_date":        &graphql.ArgumentConfig{Type: graphql.String},
			"validation_external_version":   &graphql.ArgumentConfig{Type: graphql.String},
			"validation_external_date":      &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := params.Args["mls"]
			if mls == nil {
				return nil, fmt.Errorf("You must specify the mls")
			}
			cond := map[string]interface{}{
				"id":                            helper.GetGraphParam(params, "id", ""),
				"resource_id":                   helper.GetGraphParam(params, "resource_id", ""),
				"standard_name":                 helper.GetGraphParam(params, "standard_name", ""),
				"visible_name":                  helper.GetGraphParam(params, "visible_name", ""),
				"description":                   helper.GetGraphParam(params, "description", ""),
				"key_field":                     helper.GetGraphParam(params, "key_field", ""),
				"class_count":                   helper.GetGraphParam(params, "class_count", ""),
				"class_version":                 helper.GetGraphParam(params, "class_version", ""),
				"class_date":                    helper.GetGraphParam(params, "class_date", ""),
				"object_version":                helper.GetGraphParam(params, "object_version", ""),
				"object_date":                   helper.GetGraphParam(params, "object_date", ""),
				"search_help_version":           helper.GetGraphParam(params, "search_help_version", ""),
				"search_help_date":              helper.GetGraphParam(params, "search_help_date", ""),
				"edit_mask_version":             helper.GetGraphParam(params, "edit_mask_version", ""),
				"edit_mask_date":                helper.GetGraphParam(params, "edit_mask_date", ""),
				"lookup_version":                helper.GetGraphParam(params, "lookup_version", ""),
				"lookup_date":                   helper.GetGraphParam(params, "lookup_date", ""),
				"update_help_version":           helper.GetGraphParam(params, "update_help_version", ""),
				"update_help_date":              helper.GetGraphParam(params, "update_help_date", ""),
				"validation_expression_version": helper.GetGraphParam(params, "validation_expression_version", ""),
				"validation_expression_date":    helper.GetGraphParam(params, "validation_expression_date", ""),
				"validation_lookup_version":     helper.GetGraphParam(params, "validation_lookup_version", ""),
				"validation_lookup_date":        helper.GetGraphParam(params, "validation_lookup_date", ""),
				"validation_external_version":   helper.GetGraphParam(params, "validation_external_version", ""),
				"validation_external_date":      helper.GetGraphParam(params, "validation_external_date", ""),
			}
			result := make([]map[string]interface{}, 0)
			row := orm.Model{}
			row.SetTableName(models.Resource{}.TableName())
			conn, _ := db.OpenMls(mls.(string))
			defer conn.Close()
			session := orm.NewSession(conn)
			result, err := session.Find(row, cond)
			return result, err
		},
	},
	"metadataClass": &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "MetadataClassItem",
			Fields: graphql.Fields{
				"id":             &graphql.Field{Type: graphql.String},
				"store_as":       &graphql.Field{Type: graphql.String},
				"resource_name":  &graphql.Field{Type: graphql.String},
				"resource_id":    &graphql.Field{Type: graphql.String},
				"class_name":     &graphql.Field{Type: graphql.String},
				"standard_name":  &graphql.Field{Type: graphql.String},
				"visible_name":   &graphql.Field{Type: graphql.String},
				"description":    &graphql.Field{Type: graphql.String},
				"table_version":  &graphql.Field{Type: graphql.String},
				"table_date":     &graphql.Field{Type: graphql.String},
				"mapped":         &graphql.Field{Type: graphql.Int},
				"unmapped":       &graphql.Field{Type: graphql.Int},
				"update_version": &graphql.Field{Type: graphql.String},
				"update_date":    &graphql.Field{Type: graphql.String},
			},
		})),
		Args: graphql.FieldConfigArgument{
			"mls":            &graphql.ArgumentConfig{Type: graphql.String},
			"id":             &graphql.ArgumentConfig{Type: graphql.String},
			"store_as":       &graphql.ArgumentConfig{Type: graphql.String},
			"resource_name":  &graphql.ArgumentConfig{Type: graphql.String},
			"resource_id":    &graphql.ArgumentConfig{Type: graphql.String},
			"class_name":     &graphql.ArgumentConfig{Type: graphql.String},
			"standard_name":  &graphql.ArgumentConfig{Type: graphql.String},
			"visible_name":   &graphql.ArgumentConfig{Type: graphql.String},
			"description":    &graphql.ArgumentConfig{Type: graphql.String},
			"table_version":  &graphql.ArgumentConfig{Type: graphql.String},
			"table_date":     &graphql.ArgumentConfig{Type: graphql.String},
			"update_version": &graphql.ArgumentConfig{Type: graphql.String},
			"update_date":    &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := helper.GetGraphParam(params, "mls", "")
			if mls == nil {
				return nil, fmt.Errorf("You must specify the mls")
			}
			cond := map[string]interface{}{
				"id":             helper.GetGraphParam(params, "id", ""),
				"store_as":       helper.GetGraphParam(params, "store_as", ""),
				"resource_id":    helper.GetGraphParam(params, "resource_id", ""),
				"resource_name":  helper.GetGraphParam(params, "resource_name", ""),
				"class_name":     helper.GetGraphParam(params, "class_name", ""),
				"standard_name":  helper.GetGraphParam(params, "standard_name", ""),
				"visible_name":   helper.GetGraphParam(params, "visible_name", ""),
				"description":    helper.GetGraphParam(params, "description", ""),
				"table_version":  helper.GetGraphParam(params, "table_version", ""),
				"table_date":     helper.GetGraphParam(params, "table_date", ""),
				"update_version": helper.GetGraphParam(params, "update_version", ""),
				"update_date":    helper.GetGraphParam(params, "update_date", ""),
			}
			result := make([]map[string]interface{}, 0)
			row := orm.Model{"id": 1}
			row.SetTableName(models.Class{}.TableName())
			mlsConn, err := db.OpenMls(mls.(string))
			if err != nil {
				return result, err
			}
			defer mlsConn.Close()
			session := orm.NewSession(mlsConn)
			result, err = session.Find(row, cond)
			if err != nil {
				return result, err
			}
			hostConn, err := db.Open()
			if err != nil {
				return result, err
			}
			defer hostConn.Close()

			for k, v := range result {
				class := v["class_name"].(string)
				resource := v["resource_name"].(string)
				store := v["store_as"].(string)
				if store != "" {
					mapped, unmapeed, err := metadata.CountClassMapFields(hostConn, mlsConn, resource, class)
					if err == nil {
						result[k]["mapped"] = mapped
						result[k]["unmapped"] = unmapeed
					} else {
						result[k]["mapped"] = 0
						result[k]["unmapped"] = 0
					}
				} else {
					var unmapped int64
					unmapped, err = mlsConn.
						Where("class_name = ?", class).
						Where("resource_name = ?", resource).
						Where("mapped_status != 'ALL'").
						Count(models.Table{})
					result[k]["mapped"] = 0
					result[k]["unmapped"] = unmapped
				}
			}
			return result, err
		},
	},
	"metadataTables": &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "MetadataTableItem",
			Fields: graphql.Fields{
				"id":             &graphql.Field{Type: graphql.String},
				"class_id":       &graphql.Field{Type: graphql.String},
				"class_name":     &graphql.Field{Type: graphql.String},
				"resource_id":    &graphql.Field{Type: graphql.String},
				"resource_name":  &graphql.Field{Type: graphql.String},
				"mapped":         &graphql.Field{Type: graphql.String},
				"mapped_status":  &graphql.Field{Type: graphql.String},
				"system_name":    &graphql.Field{Type: graphql.String},
				"standard_name":  &graphql.Field{Type: graphql.String},
				"long_name":      &graphql.Field{Type: graphql.String},
				"db_name":        &graphql.Field{Type: graphql.String},
				"short_name":     &graphql.Field{Type: graphql.String},
				"maximum_length": &graphql.Field{Type: graphql.String},
				"data_type":      &graphql.Field{Type: graphql.String},
				"precision":      &graphql.Field{Type: graphql.String},
				"searchable":     &graphql.Field{Type: graphql.String},
				"interpretation": &graphql.Field{Type: graphql.String},
				"alignment":      &graphql.Field{Type: graphql.String},
				"use_separator":  &graphql.Field{Type: graphql.String},
				"edit_mask_id":   &graphql.Field{Type: graphql.String},
				"lookup_name":    &graphql.Field{Type: graphql.String},
				"max_select":     &graphql.Field{Type: graphql.String},
				"units":          &graphql.Field{Type: graphql.String},
				"index":          &graphql.Field{Type: graphql.String},
				"minimum":        &graphql.Field{Type: graphql.String},
				"maximum":        &graphql.Field{Type: graphql.String},
				"default":        &graphql.Field{Type: graphql.String},
				"required":       &graphql.Field{Type: graphql.String},
				"search_help_id": &graphql.Field{Type: graphql.String},
				"unique":         &graphql.Field{Type: graphql.String},
			},
		})),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":            &graphql.ArgumentConfig{Type: graphql.String},
			"id":             &graphql.ArgumentConfig{Type: graphql.String},
			"class_id":       &graphql.ArgumentConfig{Type: graphql.String},
			"class_name":     &graphql.ArgumentConfig{Type: graphql.String},
			"resource_id":    &graphql.ArgumentConfig{Type: graphql.String},
			"resource_name":  &graphql.ArgumentConfig{Type: graphql.String},
			"mapped":         &graphql.ArgumentConfig{Type: graphql.String},
			"mapped_status":  &graphql.ArgumentConfig{Type: graphql.String},
			"system_name":    &graphql.ArgumentConfig{Type: graphql.String},
			"standard_name":  &graphql.ArgumentConfig{Type: graphql.String},
			"long_name":      &graphql.ArgumentConfig{Type: graphql.String},
			"db_name":        &graphql.ArgumentConfig{Type: graphql.String},
			"short_name":     &graphql.ArgumentConfig{Type: graphql.String},
			"maximum_length": &graphql.ArgumentConfig{Type: graphql.String},
			"data_type":      &graphql.ArgumentConfig{Type: graphql.String},
			"precision":      &graphql.ArgumentConfig{Type: graphql.String},
			"searchable":     &graphql.ArgumentConfig{Type: graphql.String},
			"interpretation": &graphql.ArgumentConfig{Type: graphql.String},
			"alignment":      &graphql.ArgumentConfig{Type: graphql.String},
			"use_separator":  &graphql.ArgumentConfig{Type: graphql.String},
			"edit_mask_id":   &graphql.ArgumentConfig{Type: graphql.String},
			"lookup_name":    &graphql.ArgumentConfig{Type: graphql.String},
			"max_select":     &graphql.ArgumentConfig{Type: graphql.String},
			"units":          &graphql.ArgumentConfig{Type: graphql.String},
			"index":          &graphql.ArgumentConfig{Type: graphql.String},
			"minimum":        &graphql.ArgumentConfig{Type: graphql.String},
			"maximum":        &graphql.ArgumentConfig{Type: graphql.String},
			"default":        &graphql.ArgumentConfig{Type: graphql.String},
			"required":       &graphql.ArgumentConfig{Type: graphql.String},
			"search_help_id": &graphql.ArgumentConfig{Type: graphql.String},
			"unique":         &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := params.Args["mls"]
			if mls == nil {
				return nil, fmt.Errorf("You must specify the mls")
			}
			cond := map[string]interface{}{
				"id":             helper.GetGraphParam(params, "id", ""),
				"class_id":       helper.GetGraphParam(params, "class_id", ""),
				"class_name":     helper.GetGraphParam(params, "class_name", ""),
				"resource_id":    helper.GetGraphParam(params, "resource_id", ""),
				"resource_name":  helper.GetGraphParam(params, "resource_name", ""),
				"mapped":         helper.GetGraphParam(params, "mapped", ""),
				"mapped_status":  helper.GetGraphParam(params, "mapped", ""),
				"system_name":    helper.GetGraphParam(params, "system_name", ""),
				"standard_name":  helper.GetGraphParam(params, "standard_name", ""),
				"long_name":      helper.GetGraphParam(params, "long_name", ""),
				"db_name":        helper.GetGraphParam(params, "db_name", ""),
				"short_name":     helper.GetGraphParam(params, "short_name", ""),
				"maximum_length": helper.GetGraphParam(params, "maximum_length", ""),
				"data_type":      helper.GetGraphParam(params, "data_type", ""),
				"precision":      helper.GetGraphParam(params, "precision", ""),
				"searchable":     helper.GetGraphParam(params, "searchable", ""),
				"interpretation": helper.GetGraphParam(params, "interpretation", ""),
				"alignment":      helper.GetGraphParam(params, "alignment", ""),
				"use_separator":  helper.GetGraphParam(params, "use_separator", ""),
				"edit_mask_id":   helper.GetGraphParam(params, "edit_mask_id", ""),
				"lookup_name":    helper.GetGraphParam(params, "lookup_name", ""),
				"max_select":     helper.GetGraphParam(params, "max_select", ""),
				"units":          helper.GetGraphParam(params, "units", ""),
				"index":          helper.GetGraphParam(params, "index", ""),
				"minimum":        helper.GetGraphParam(params, "minimum", ""),
				"maximum":        helper.GetGraphParam(params, "maximum", ""),
				"default":        helper.GetGraphParam(params, "default", ""),
				"required":       helper.GetGraphParam(params, "required", ""),
				"search_help_id": helper.GetGraphParam(params, "search_help_id", ""),
				"unique":         helper.GetGraphParam(params, "unique", ""),
			}
			result := make([]map[string]interface{}, 0)
			row := orm.Model{}
			row.SetTableName(models.Table{}.TableName())
			conn, _ := db.OpenMls(mls.(string))
			defer conn.Close()
			session := orm.NewSession(conn)
			result, err := session.Find(row, cond)
			return result, err
		},
	},
}
