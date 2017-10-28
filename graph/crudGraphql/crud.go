package crudGraphql

import (
	"fmt"
	"strings"

	"bitbucket.org/mlsdatatools/retsetl/api/crud"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
	"github.com/graphql-go/graphql"
)

//NewCrud Create the necessary structure to run a crud on any table under the api graphql
func NewCrud(conn *xorm.Engine, model string) (GetData graphql.Fields, SetData graphql.Fields) {
	session := orm.NewSession(conn)
	tmp := orm.Model{}
	tmp.SetTableName(model)
	fields := orm.GetFieldsDB(session, tmp)
	GetData = newGetData(model, fields)
	SetData = newSetData(model, fields)
	return
}

func getModelCamel(model string) string {
	t := strings.Split(model, "_")
	modelCamel := ""
	if len(t) > 1 {
		modelCamel = t[0]
		for k := 0; k < len(t); k++ {
			modelCamel = modelCamel + strings.Title(t[k])
		}
	} else {
		modelCamel = model
	}
	return modelCamel
}

func newGetData(model string, fields []string) graphql.Fields {
	modelCamel := getModelCamel(model)
	query := fmt.Sprintf("%vQuery", modelCamel)
	get := graphql.Fields{}
	get[query] = newGetDataQuery(model, fields)
	return get
}

func newGetDataQuery(model string, fields []string) *graphql.Field {
	modelCamel := getModelCamel(model)
	Fields := graphql.Fields{}
	FieldsConfigArgument := graphql.FieldConfigArgument{}
	for _, v := range fields {
		Fields[v] = &graphql.Field{Type: graphql.String}
		FieldsConfigArgument[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	return &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name:   fmt.Sprintf("%vItem", modelCamel),
			Fields: Fields,
		})),
		Description: "",
		Args:        FieldsConfigArgument,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			var data []map[string]interface{}
			row := orm.Model{}
			row.SetTableName(model)
			cond := builder.Eq{}
			for k, v := range params.Args {
				cond[k] = v.(string)
			}
			data, _, err = crud.Search(row, cond)
			return data, err
		},
	}
}

func newSetDataCreate(model string, fields []string) *graphql.Field {
	modelCamel := getModelCamel(model)
	Fields := graphql.Fields{}
	Args := graphql.FieldConfigArgument{}
	name := fmt.Sprintf("%vItemIsert", modelCamel)
	for _, v := range fields {
		Fields[v] = &graphql.Field{Type: graphql.String}
		Args[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	delete(Args, "id")
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   name,
			Fields: Fields,
		}),
		Description: "",
		Args:        Args,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			row := orm.Model{}
			for k, v := range params.Args {
				row[k] = fmt.Sprintf("%v", v)
			}
			row.SetTableName(model)
			_, err = crud.Create(row)
			response := map[string]interface{}{}
			for k, v := range row {
				response[k] = v
			}
			return response, err
		},
	}
}

func newSetDataUpdate(model string, fields []string) *graphql.Field {
	modelCamel := getModelCamel(model)
	Fields := graphql.Fields{}
	Args := graphql.FieldConfigArgument{}
	for _, v := range fields {
		Fields[v] = &graphql.Field{Type: graphql.String}
		Args[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   fmt.Sprintf("%vItemUpdate", modelCamel),
			Fields: Fields,
		}),
		Description: "",
		Args:        Args,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			row := orm.Model{}
			for k, v := range params.Args {
				row[k] = fmt.Sprintf("%v", v)
			}
			row.SetTableName(model)
			_, err = crud.Update(row)
			response := map[string]interface{}{}
			for k, v := range row {
				response[k] = v
			}
			return response, err
		},
	}
}

func newSetDataDelete(model string, fields []string) *graphql.Field {
	modelCamel := getModelCamel(model)
	Fields := graphql.Fields{}
	Args := graphql.FieldConfigArgument{}
	for _, v := range fields {
		Fields[v] = &graphql.Field{Type: graphql.String}
		Args[v] = &graphql.ArgumentConfig{Type: graphql.String}
	}
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   fmt.Sprintf("%vItemDelete", modelCamel),
			Fields: Fields,
		}),
		Description: "",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var err error
			row := orm.Model{}
			if params.Args["id"] == nil {
				return nil, fmt.Errorf("You must specify the id")
			}
			row["id"] = params.Args["id"].(string)
			row.SetTableName(model)
			_, err = crud.Delete(row)
			response := map[string]interface{}{}
			for k, v := range row {
				response[k] = v
			}
			return response, err
		},
	}
}

func newSetData(model string, fields []string) graphql.Fields {
	modelCamel := getModelCamel(model)
	create := fmt.Sprintf("%vCreate", modelCamel)
	update := fmt.Sprintf("%vUpdate", modelCamel)
	delete := fmt.Sprintf("%vDelete", modelCamel)
	set := graphql.Fields{}
	set[create] = newSetDataCreate(model, fields)
	set[update] = newSetDataUpdate(model, fields)
	set[delete] = newSetDataDelete(model, fields)
	return set
}
