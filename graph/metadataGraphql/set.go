package metadataGraphql

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/api/rets"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"github.com/graphql-go/graphql"
)

//SetData Contains the anchor points graphql on which to modify the data
var SetData = graphql.Fields{
	"metadataClassCompare": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "metadataClassCompareSuccess",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":      &graphql.ArgumentConfig{Type: graphql.String},
			"class_id": &graphql.ArgumentConfig{Type: graphql.String},
			"resource": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := params.Args["mls"]
			classID := params.Args["class_id"]
			resource := params.Args["resource"]
			if mls == nil {
				return nil, fmt.Errorf("You must specify the mls")
			}
			if classID == nil {
				return nil, fmt.Errorf("You must specify the class_id")
			}
			if resource == nil {
				return nil, fmt.Errorf("You must specify the resource")
			}
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			mlsConn, err := db.OpenMls(mls.(string))
			if err != nil {
				return nil, err
			}
			defer mlsConn.Close()
			_mls, _classID, _resource := mls.(string), classID.(string), resource.(string)
			err = rets.CompareResourceClass(hostConn, mlsConn, _mls, _classID, _resource)
			return map[string]interface{}{"message": "OK"}, err
		},
	},
	"metadataGetResources": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "metadataImportSuccess",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := params.Args["mls"]
			if mls == nil {
				return nil, fmt.Errorf("You must specify the id")
			}
			conn, err := db.OpenMls(mls.(string))
			defer conn.Close()
			if err != nil {
				return map[string]string{"message": "Falure"}, err
			}
			result := rets.GetResources(conn, mls.(string))
			return map[string]interface{}{"message": "OK"}, result
		},
	},
	"metadataMapResource": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "metadataMapResource",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":      &graphql.ArgumentConfig{Type: graphql.String},
			"class_id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := params.Args["mls"]
			classID := params.Args["class_id"]
			if mls == nil {
				return nil, fmt.Errorf("You must specify the mls")
			}
			if classID == nil {
				return nil, fmt.Errorf("You must specify the class_id")
			}
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			mlsConn, err := db.OpenMls(mls.(string))
			if err != nil {
				return nil, err
			}
			defer mlsConn.Close()
			mapResource := rets.NewMapResource(hostConn, mlsConn, classID.(string))
			err = mapResource.Map()
			return map[string]interface{}{"message": "OK"}, err
		},
	},
	"metadataClassStoreAs": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "metadataClassStoreAs",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":      &graphql.ArgumentConfig{Type: graphql.String},
			"store_as": &graphql.ArgumentConfig{Type: graphql.String},
			"class_id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := helper.GetGraphParam(params, "mls", "")
			storeAs := helper.GetGraphParam(params, "store_as", "")
			classID := helper.GetGraphParam(params, "class_id", "")
			if mls == "" {
				return nil, fmt.Errorf("You must specify the mls")
			}
			if classID == "" {
				return nil, fmt.Errorf("You must specify the class_id")
			}
			conn, err := db.OpenMls(mls.(string))
			defer conn.Close()
			rets.ClassStoreAs(conn, classID.(string), storeAs.(string))
			return map[string]interface{}{"message": "OK"}, err
		},
	},
}
