package dictionaryGraphql

import (
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/mlsdatatools/retsetl/models"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"github.com/graphql-go/graphql"
)

//saveSynonym Stores the synonym of a field in the system
func saveSynonym() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "SaveSynonymSuccess",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"mls":          &graphql.ArgumentConfig{Type: graphql.String},
			"Synonym":      &graphql.ArgumentConfig{Type: graphql.String},
			"FieldID":      &graphql.ArgumentConfig{Type: graphql.String},
			"resource":     &graphql.ArgumentConfig{Type: graphql.String},
			"StandardName": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			mls := helper.GetGraphParam(params, "mls", "").(string)
			resource := helper.GetGraphParam(params, "resource", "").(string)
			Synonym := helper.GetGraphParam(params, "Synonym", "").(string)
			StandardName := helper.GetGraphParam(params, "StandardName", "").(string)
			if mls == "" {
				return nil, fmt.Errorf("8000: mls not found")
			}
			if StandardName == "" {
				return nil, fmt.Errorf("8000: StandardName not found")
			}
			mlsConn, err := db.OpenMls(mls)
			if err != nil {
				return nil, fmt.Errorf("8001: %v", err)
			}
			defer mlsConn.Close()
			hostConn, err := db.Open()
			if err != nil {
				return nil, fmt.Errorf("8011: %v", err)
			}
			defer hostConn.Close()

			row := models.DictionaryResource{
				Synonym:     Synonym,
				FieldStatus: "ALL",
				Table:       resource,
			}
			_, err = hostConn.
				Cols("Synonym", "FieldStatus").
				Where("StandardName = ?", StandardName).
				Update(row)
			if err != nil {
				return nil, fmt.Errorf("8004 %v", err)
			}

			table := models.Table{}
			_, err = mlsConn.
				Where("mapped = ?", StandardName).
				Where("resource_name = ?", resource).
				Get(&table)
			if err != nil {
				return nil, fmt.Errorf("8005 %v", err)
			}
			if table.ID != "" {
				table.Mapped = "none"
				table.MappedStatus = "none"
				_, err = mlsConn.Id(table.ID).Cols("mapped", "mapped_status").Update(table)
				if err != nil {
					return nil, fmt.Errorf("8006 %v", err)
				}
			}

			table = models.Table{}
			_Synonym := strings.Replace(Synonym, "'", "\\'", -1)

			valid, _ := regexp.Compile("[a-zA-Z]")
			if valid.Match([]byte(_Synonym)) {
				pieces := "'" + strings.Join(strings.Split(_Synonym, ","), "', '") + "'"
				cond := fmt.Sprintf(" in (%v)", pieces)
				_, err = mlsConn.
					Where(fmt.Sprintf("%v %v", "short_name", cond)).
					Or(fmt.Sprintf("%v %v", "long_name", cond)).
					Or(fmt.Sprintf("%v %v", "standard_name", cond)).
					Get(&table)
				if err != nil {
					return nil, fmt.Errorf("8005 %v", err)
				}
				if table.ID != "" {
					table.Mapped = StandardName
					table.MappedStatus = "ALL"
					_, err = mlsConn.Id(table.ID).Cols("mapped", "mapped_status").Update(table)
					if err != nil {
						return nil, fmt.Errorf("8006 %v", err)
					}
				}
			}
			return map[string]interface{}{"message": "OK"}, nil
		},
	}
}
