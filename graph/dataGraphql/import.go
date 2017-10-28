package dataGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/helper"

	"bitbucket.org/mlsdatatools/retsetl/api/data"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/graphql-go/graphql"
)

//SetData Contains the anchor points graphql on which to modify the data
var SetData = graphql.Fields{
	"dataImport": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "dataImportSuccess",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":       &graphql.ArgumentConfig{Type: graphql.String},
			"class_id":  &graphql.ArgumentConfig{Type: graphql.String},
			"client_id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := helper.GetGraphParam(params, "mls", "").(string)
			classID := helper.GetGraphParam(params, "class_id", "").(string)
			clientID := helper.GetGraphParam(params, "client_id", "").(string)
			if mls == "" {
				return nil, fmt.Errorf("You must specify the mls")
			}
			if classID == "" {
				return nil, fmt.Errorf("class_id required")
			}
			if clientID == "" {
				return nil, fmt.Errorf("client_id required")
			}
			hostConn, err := db.Open()
			defer hostConn.Close()
			mlsConn, err := db.OpenMls(mls)
			if err != nil {
				return nil, err
			}
			defer mlsConn.Close()

			row := models.Mls{}
			mlsConn.
				Where("id = ?", mls).
				Find(row)
			remoteConn, err := db.OpenMlsConn(hostConn, mls)
			if err != nil {
				return nil, err
			}
			defer mlsConn.Close()
			err = data.ImportData(hostConn, mlsConn, remoteConn, clientID, classID)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"message": "OK"}, err
		},
	},
}
