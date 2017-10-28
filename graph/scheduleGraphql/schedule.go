package scheduleGraphql

import (
	"strings"

	"fmt"

	"regexp"

	apidata "bitbucket.org/mlsdatatools/retsetl/api/data"
	"bitbucket.org/mlsdatatools/retsetl/api/schedule"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/graphql-go/graphql"
)

var data = graphql.NewObject(graphql.ObjectConfig{
	Name: "Schedule",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"sunday":    &graphql.Field{Type: graphql.Boolean},
		"monday":    &graphql.Field{Type: graphql.Boolean},
		"tuesday":   &graphql.Field{Type: graphql.Boolean},
		"wednesday": &graphql.Field{Type: graphql.Boolean},
		"thursday":  &graphql.Field{Type: graphql.Boolean},
		"friday":    &graphql.Field{Type: graphql.Boolean},
		"saturday":  &graphql.Field{Type: graphql.Boolean},
		"frequency": &graphql.Field{Type: graphql.Int},
		"hour":      &graphql.Field{Type: graphql.String},
		"classes":   &graphql.Field{Type: graphql.String},
		"status":    &graphql.Field{Type: graphql.String},
	},
})

//GetData Contains the anchor points graphql on which to consult the data
var GetData = graphql.Fields{
	"scheduleRun": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "scheduleRunResult",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			id := helper.GetGraphParam(params, "id", "")
			sql := fmt.Sprintf(`SELECT * FROM schedule WHERE id = '%v'`, id)
			if result, err := hostConn.QueryString(sql); err == nil {
				row := models.Schedule{Status: "Importing"}
				_, err = hostConn.ID(id).Cols("status").Update(row)
				if err != nil {
					return nil, err
				}
				go func() {
					for _, v := range result {
						classes := strings.Split(v["classes"], ",")
						for _, class := range classes {
							if mlsConn, err := db.OpenMls(v["mls"]); err == nil {
								remoteConn, err := db.OpenMlsConn(hostConn, v["mls"])
								defer mlsConn.Close()
								if err != nil {
									continue
								}
								if err == nil {
									apidata.ImportData(hostConn, mlsConn, remoteConn, "", class)
								}
							}
						}
					}
				}()
			} else {
				fmt.Println(err)
			}
			return "OK", nil
		},
	},
	"scheduleQuery": &graphql.Field{
		Type:        graphql.NewList(data),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			rows := make([]models.Schedule, 0)
			err = hostConn.Find(&rows)
			return rows, err
		},
	},
}

//SetData Contains the anchor points graphql on which to modify the data
var SetData = graphql.Fields{
	"scheduleCreate": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "scheduleCreateResult",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Args: graphql.FieldConfigArgument{
			"mls":       &graphql.ArgumentConfig{Type: graphql.String},
			"sunday":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			"monday":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			"tuesday":   &graphql.ArgumentConfig{Type: graphql.Boolean},
			"wednesday": &graphql.ArgumentConfig{Type: graphql.Boolean},
			"thursday":  &graphql.ArgumentConfig{Type: graphql.Boolean},
			"friday":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			"saturday":  &graphql.ArgumentConfig{Type: graphql.Boolean},
			"frequency": &graphql.ArgumentConfig{Type: graphql.Int},
			"hour":      &graphql.ArgumentConfig{Type: graphql.String},
			"classes":   &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			mls := helper.GetGraphParam(params, "mls", "")
			sunday := helper.GetGraphParam(params, "sunday", false)
			monday := helper.GetGraphParam(params, "monday", false)
			tuesday := helper.GetGraphParam(params, "tuesday", false)
			wednesday := helper.GetGraphParam(params, "wednesday", false)
			thursday := helper.GetGraphParam(params, "thursday", false)
			friday := helper.GetGraphParam(params, "friday", false)
			saturday := helper.GetGraphParam(params, "saturday", false)
			frequency := helper.GetGraphParam(params, "frequency", 1)
			hour := helper.GetGraphParam(params, "hour", "00:00:00")
			classes := helper.GetGraphParam(params, "classes", "")
			p, _ := regexp.Compile("([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]")
			if !p.MatchString(hour.(string)) {
				return nil, fmt.Errorf("hour no valid")
			}
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			fmt.Println(classes)
			err = schedule.CreateSchedule(
				hostConn,
				mls.(string),
				sunday.(bool),
				monday.(bool),
				tuesday.(bool),
				wednesday.(bool),
				thursday.(bool),
				friday.(bool),
				saturday.(bool),
				frequency.(int),
				hour.(string),
				classes.(string),
			)
			return "OK", err
		},
	},
	"scheduleUpdate": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "scheduleUpdateResult",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Args: graphql.FieldConfigArgument{
			"id":        &graphql.ArgumentConfig{Type: graphql.String},
			"mls":       &graphql.ArgumentConfig{Type: graphql.String},
			"sunday":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			"monday":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			"tuesday":   &graphql.ArgumentConfig{Type: graphql.Boolean},
			"wednesday": &graphql.ArgumentConfig{Type: graphql.Boolean},
			"thursday":  &graphql.ArgumentConfig{Type: graphql.Boolean},
			"friday":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			"saturday":  &graphql.ArgumentConfig{Type: graphql.Boolean},
			"frequency": &graphql.ArgumentConfig{Type: graphql.Int},
			"hour":      &graphql.ArgumentConfig{Type: graphql.String},
			"classes":   &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ID := helper.GetGraphParam(params, "id", false)
			mls := helper.GetGraphParam(params, "mls", "")
			sunday := helper.GetGraphParam(params, "sunday", false)
			monday := helper.GetGraphParam(params, "monday", false)
			tuesday := helper.GetGraphParam(params, "tuesday", false)
			wednesday := helper.GetGraphParam(params, "wednesday", false)
			thursday := helper.GetGraphParam(params, "thursday", false)
			friday := helper.GetGraphParam(params, "friday", false)
			saturday := helper.GetGraphParam(params, "saturday", false)
			frequency := helper.GetGraphParam(params, "frequency", 1)
			hour := helper.GetGraphParam(params, "hour", "00:00:00")
			classes := helper.GetGraphParam(params, "classes", "")
			p, _ := regexp.Compile("([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]")
			if !p.MatchString(hour.(string)) {
				return nil, fmt.Errorf("hour no valid")
			}
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			err = schedule.UpdateSchedule(
				hostConn,
				ID.(string),
				mls.(string),
				sunday.(bool),
				monday.(bool),
				tuesday.(bool),
				wednesday.(bool),
				thursday.(bool),
				friday.(bool),
				saturday.(bool),
				frequency.(int),
				hour.(string),
				classes.(string),
			)
			return "OK", err
		},
	},
	"scheduleDelete": &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "scheduleDeleteResult",
			Fields: graphql.Fields{
				"message": &graphql.Field{Type: graphql.String},
			},
		}),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ID := helper.GetGraphParam(params, "id", false)
			hour := helper.GetGraphParam(params, "", "00:00:00")
			p, _ := regexp.Compile("([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]")
			if !p.MatchString(hour.(string)) {
				return nil, fmt.Errorf("hour no valid")
			}
			hostConn, err := db.Open()
			if err != nil {
				return nil, err
			}
			defer hostConn.Close()
			err = schedule.DeleteSchedule(hostConn, ID.(string))
			return "OK", err
		},
	},
}
