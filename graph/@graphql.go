package graph

import (
	"bitbucket.org/mlsdatatools/retsetl/graph/dataGraphql"
	"bitbucket.org/mlsdatatools/retsetl/graph/dictionaryGraphql"
	"bitbucket.org/mlsdatatools/retsetl/graph/logsGraphql"
	"bitbucket.org/mlsdatatools/retsetl/graph/metadataGraphql"
	"bitbucket.org/mlsdatatools/retsetl/graph/mlsGraphql"
	"bitbucket.org/mlsdatatools/retsetl/graph/scheduleGraphql"
	"github.com/graphql-go/graphql"
)

var schema graphql.Schema

//ExecuteQuery Run graphql queries
func ExecuteQuery(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}

func concat(store graphql.Fields, append graphql.Fields) graphql.Fields {
	for i, v := range append {
		store[i] = v
	}
	return store
}

func init() {
	var query = make(graphql.Fields, 0)
	var mutation = make(graphql.Fields, 0)
	query = concat(query, mlsGraphql.GetData)
	mutation = concat(mutation, mlsGraphql.SetData)

	query = concat(query, metadataGraphql.GetData)
	mutation = concat(mutation, metadataGraphql.SetData)

	mutation = concat(mutation, dataGraphql.SetData)

	query = concat(query, dictionaryGraphql.GetData)
	mutation = concat(mutation, dictionaryGraphql.SetData)

	query = concat(query, logsGraphql.GetData)

	query = concat(query, scheduleGraphql.GetData)
	mutation = concat(mutation, scheduleGraphql.SetData)

	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: query,
	})

	var rootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: mutation,
	})

	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}
