package mlsGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/config"

	"bitbucket.org/mlsdatatools/retsetl/api/mls"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/graphql-go/graphql"
	"gopkg.in/mgo.v2/bson"
)

var data = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mls",
	Fields: graphql.Fields{
		"id":             &graphql.Field{Type: graphql.String},
		"name":           &graphql.Field{Type: graphql.String},
		"system_id":      &graphql.Field{Type: graphql.String},
		"url":            &graphql.Field{Type: graphql.String},
		"username":       &graphql.Field{Type: graphql.String},
		"useragent_name": &graphql.Field{Type: graphql.String},
		"version":        &graphql.Field{Type: graphql.String},

		"server_bd":   &graphql.Field{Type: graphql.String},
		"username_bd": &graphql.Field{Type: graphql.String},
		"password_bd": &graphql.Field{Type: graphql.String},
		"name_bd":     &graphql.Field{Type: graphql.String},
	},
})

//GetData Contains the anchor points graphql on which to consult the data
var GetData = graphql.Fields{
	"mlsQuery": &graphql.Field{
		Type:        graphql.NewList(data),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			query := map[string]interface{}{}
			result := make([]models.Mls, 0)
			if params.Args["id"] != nil {
				query["id"] = params.Args["id"].(string)
				_, err = mls.Search(&result, params.Args["id"].(string))
			} else {
				_, err = mls.Search(&result)
			}
			return result, err
		},
	},
}

//SetData Contains the anchor points graphql on which to modify the data
var SetData = graphql.Fields{
	"mlsCreate": &graphql.Field{
		Type: responseCreateUpdate,
		Args: graphql.FieldConfigArgument{
			"url":                &graphql.ArgumentConfig{Type: graphql.String},
			"username":           &graphql.ArgumentConfig{Type: graphql.String},
			"password":           &graphql.ArgumentConfig{Type: graphql.String},
			"useragent_name":     &graphql.ArgumentConfig{Type: graphql.String},
			"useragent_password": &graphql.ArgumentConfig{Type: graphql.String},
			"version":            &graphql.ArgumentConfig{Type: graphql.String},

			"server_bd":   &graphql.ArgumentConfig{Type: graphql.String},
			"username_bd": &graphql.ArgumentConfig{Type: graphql.String},
			"password_bd": &graphql.ArgumentConfig{Type: graphql.String},
			"name_bd":     &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			row := models.Mls{}

			row.ID = bson.NewObjectId().Hex()
			row.URL = helper.GetGraphParam(params, "url", "").(string)
			row.Username = helper.GetGraphParam(params, "username", "").(string)
			row.Password = helper.GetGraphParam(params, "password", "").(string)
			row.UseragentName = helper.GetGraphParam(params, "useragent_name", "").(string)
			row.UseragentPassword = helper.GetGraphParam(params, "useragent_password", "").(string)
			row.VersionRets = helper.GetGraphParam(params, "version", "").(string)

			serverBd := helper.GetGraphParam(params, "server_bd", "").(string)
			usernameBd := helper.GetGraphParam(params, "username_bd", "").(string)
			passwordBd := helper.GetGraphParam(params, "password_bd", "").(string)
			nameBd := helper.GetGraphParam(params, "name_bd", "").(string)

			if serverBd+usernameBd+passwordBd+nameBd == "" {
				serverBd = "199.168.136.134"
				usernameBd = "mlsvowuser"
				passwordBd = "0112358"
				nameBd = row.ID
				conn, err := db.OpenBD(
					config.DB_USER,
					config.DB_PWD,
					config.DB_HOST,
					config.DB_PORT,
					config.DB_NAME,
					"postgres",
				)
				if err != nil {
					return nil, err
				}
				sql := `CREATE DATABASE "` + nameBd + `";`
				if _, err := conn.Exec(sql); err != nil {
					fmt.Println(sql, err)
					return nil, err
				}
			}

			row.ServerBd = serverBd
			row.UsernameBd = usernameBd
			row.PasswordBd = passwordBd
			row.NameBd = nameBd

			_, err := mls.Create(&row)
			response := structCreateUpdate{
				Message: "OK",
				Record:  row,
			}
			return response, err
		},
	},
	"mlsUpdate": &graphql.Field{
		Type: responseCreateUpdate,
		Args: graphql.FieldConfigArgument{
			"id":                 &graphql.ArgumentConfig{Type: graphql.String},
			"url":                &graphql.ArgumentConfig{Type: graphql.String},
			"username":           &graphql.ArgumentConfig{Type: graphql.String},
			"password":           &graphql.ArgumentConfig{Type: graphql.String},
			"useragent_name":     &graphql.ArgumentConfig{Type: graphql.String},
			"useragent_password": &graphql.ArgumentConfig{Type: graphql.String},
			"version":            &graphql.ArgumentConfig{Type: graphql.String},

			"server_bd":   &graphql.ArgumentConfig{Type: graphql.String},
			"username_bd": &graphql.ArgumentConfig{Type: graphql.String},
			"password_bd": &graphql.ArgumentConfig{Type: graphql.String},
			"name_bd":     &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			id := helper.GetGraphParam(params, "id", "").(string)
			row := models.Mls{}
			row.ID = id

			row.URL = helper.GetGraphParam(params, "url", "").(string)
			row.Username = helper.GetGraphParam(params, "username", "").(string)
			row.Password = helper.GetGraphParam(params, "password", "").(string)
			row.UseragentName = helper.GetGraphParam(params, "useragent_name", "").(string)
			row.UseragentPassword = helper.GetGraphParam(params, "useragent_password", "").(string)
			row.VersionRets = helper.GetGraphParam(params, "version", "").(string)

			row.ServerBd = helper.GetGraphParam(params, "server_bd", "").(string)
			row.UsernameBd = helper.GetGraphParam(params, "username_bd", "").(string)
			row.PasswordBd = helper.GetGraphParam(params, "password_bd", "").(string)
			row.NameBd = helper.GetGraphParam(params, "name_bd", "").(string)
			_, err := mls.Update(&row)
			response := structCreateUpdate{
				Message: "OK",
				Record:  row,
			}
			return response, err
		},
	},
	"mlsDelete": &graphql.Field{
		Type: responseCreateUpdate,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			row := models.Mls{}
			id := helper.GetGraphParam(params, "id", "").(string)
			_, err := mls.Delete(id, &row)
			response := structCreateUpdate{
				Message: "OK",
				Record:  row,
			}
			return response, err
		},
	},
}

var responseCreateUpdate = graphql.NewObject(graphql.ObjectConfig{
	Name: "result",
	Fields: graphql.Fields{
		"message": &graphql.Field{Type: graphql.String},
		"record":  &graphql.Field{Type: data},
	},
})

type structCreateUpdate struct {
	Message string     `json:"message"`
	Record  models.Mls `json:"record"`
}
