package orm

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/config"
)

//Model Represents a record of the unknown table
type Model map[string]interface{}

//GetFieldsDB Get fields from a table in bd
func GetFieldsDB(session *Session, model Model) []string {
	switch config.DRIVER {
	case "sqlite3":
		tmp := "pragma table_info('%v');"
		sql := fmt.Sprintf(tmp, model.TableName())
		result, _ := session.engine.Query(sql)
		response := make([]string, 0)
		for _, row := range result {
			response = append(response, string(row["name"]))
		}
		return response
	case "postgres":
		tmp := "SELECT column_name FROM information_schema.columns WHERE table_name = '%v'"
		sql := fmt.Sprintf(tmp, model.TableName())
		result, _ := session.engine.Query(sql)
		response := make([]string, 0)
		for _, row := range result {
			response = append(response, string(row["column_name"]))
		}
		return response
	}
	result := make([]string, 0)
	return result
}

//SetTableName Sets the name of the table to be used.
func (el Model) SetTableName(str string) {
	el["tablename"] = str
}

//TableName This function returns the name of the table in the database that is used to store the data
func (el Model) TableName() string {
	if value, ok := el["tablename"]; ok {
		return fmt.Sprintf("%v", value)
	}
	return ""
}

//Cols Gets the names of the columns in the model, this value is calculated from the assigned data
func (el Model) Cols() []string {
	fields := make([]string, 0)
	for k := range el {
		if k != "tablename" {
			fields = append(fields, k)
		}
	}
	return fields
}

//Values It obtains an array of the values that the model owns,
//the order of the values corresponds with the order of the columns
//obtained with the method Cols
func (el Model) Values() []interface{} {
	values := make([]interface{}, 0)
	for k, v := range el {
		if k != "tablename" {
			values = append(values, v)
		}
	}
	return values
}
