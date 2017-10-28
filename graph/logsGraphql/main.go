package logsGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/helper"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"

	"github.com/graphql-go/graphql"
)

//GetData Contains the anchor points graphql on which to consult the data
var GetData = graphql.Fields{
	"logsAutoMap": &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "logsAutoMapItem",
			Fields: graphql.Fields{
				"id":            &graphql.Field{Type: graphql.String},
				"log_id":        &graphql.Field{Type: graphql.String},
				"synonymous":    &graphql.Field{Type: graphql.String},
				"max_length":    &graphql.Field{Type: graphql.String},
				"type":          &graphql.Field{Type: graphql.String},
				"standard_name": &graphql.Field{Type: graphql.String},
				"status":        &graphql.Field{Type: graphql.String},
				"system_name":   &graphql.Field{Type: graphql.String},
				"resource_name": &graphql.Field{Type: graphql.String},
				"field_id":      &graphql.Field{Type: graphql.String},
				"class_id":      &graphql.Field{Type: graphql.String},
				"resource_id":   &graphql.Field{Type: graphql.String},
				"compare":       &graphql.Field{Type: graphql.String},
			},
		})),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"class_id": &graphql.ArgumentConfig{Type: graphql.String},
			"mls":      &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			classID, mls := params.Args["class_id"], params.Args["mls"]
			if classID == nil || mls == nil {
				return nil, fmt.Errorf("class_id or mls not found")
			}
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			mlsCon, err := db.OpenMls(mls.(string))
			if err != nil {
				return nil, err
			}
			class := models.Class{}
			_, err = mlsCon.
				Where("id = ?", classID).
				Get(&class)
			log := models.Log{}
			_, err = hostConn.
				Where("resource_name = ?", class.StoreAs).
				Where("class_id = ?", classID).
				Where("mls = ?", mls).
				OrderBy("date desc").
				Get(&log)
			details := make([]models.LogDetails, 0)
			err = hostConn.
				Where("log_id = ?", log.ID).
				Find(&details)
			return details, err
		},
	},
	"logsImport": &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "logsImportItem",
			Fields: graphql.Fields{
				"id":          &graphql.Field{Type: graphql.String},
				"mls":         &graphql.Field{Type: graphql.String},
				"resource":    &graphql.Field{Type: graphql.String},
				"resource_id": &graphql.Field{Type: graphql.String},
				"class":       &graphql.Field{Type: graphql.String},
				"class_id":    &graphql.Field{Type: graphql.String},
				"date":        &graphql.Field{Type: graphql.String},
				"imported":    &graphql.Field{Type: graphql.String},
				"status":      &graphql.Field{Type: graphql.String},
			},
		})),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"id":          &graphql.ArgumentConfig{Type: graphql.String},
			"mls":         &graphql.ArgumentConfig{Type: graphql.String},
			"resource":    &graphql.ArgumentConfig{Type: graphql.String},
			"resource_id": &graphql.ArgumentConfig{Type: graphql.String},
			"class":       &graphql.ArgumentConfig{Type: graphql.String},
			"class_id":    &graphql.ArgumentConfig{Type: graphql.String},
			"date":        &graphql.ArgumentConfig{Type: graphql.String},
			"imported":    &graphql.ArgumentConfig{Type: graphql.String},
			"status":      &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := helper.GetGraphParam(params, "mls", "").(string)
			if mls == "" {
				return nil, fmt.Errorf("mls not found")
			}
			result := make([]models.ImportLog, 0)
			hostConn, err := db.Open()
			if err != nil {
				return result, err
			}
			session := hostConn.NewSession()
			for k, v := range params.Args {
				session.Where(k+" = ?", v)
			}
			err = session.Find(&result)
			return result, err
		},
	},
	"logsImportDetails": &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "logsImportDetailsItem",
			Fields: graphql.Fields{
				"id":            &graphql.Field{Type: graphql.String},
				"import_log_id": &graphql.Field{Type: graphql.String},
				"error":         &graphql.Field{Type: graphql.String},
				"data":          &graphql.Field{Type: graphql.String},
			},
		})),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"import_log_id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			importLogID := helper.GetGraphParam(params, "import_log_id", "").(string)
			result := make([]models.ImportLogDetails, 0)
			if importLogID == "" {
				return result, fmt.Errorf("import_log_id not found")
			}
			hostConn, err := db.Open()
			if err != nil {
				return result, err
			}
			err = hostConn.
				Where("import_log_id = ?", importLogID).
				Find(&result)
			return result, err
		},
	},
}
