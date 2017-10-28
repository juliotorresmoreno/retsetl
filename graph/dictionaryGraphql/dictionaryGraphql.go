package dictionaryGraphql

import (
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/graph/crudGraphql"
	"github.com/graphql-go/graphql"
)

//GetData Apis exposed in the query
var GetData = graphql.Fields{
	"resourceQuery": query(), //Query the fields of any resource in the system
}

//SetData Apis exposed in the mutations
var SetData = graphql.Fields{
	"resourceSaveSynonym": saveSynonym(), //Stores the synonym of a field in the system
	"resourceCreate":      create(),      //Create the fields of any resource in the system
	"resourceUpdate":      update(),      //Update the fields of any resource in the system
	"resourceDelete":      remove(),      //Delete the fields of any resource in the system
}

func init() {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		dictionaryGetData, _ := crudGraphql.NewCrud(conn, "LookupFieldsAndValues")
		GetData = concat(GetData, dictionaryGetData)
	}
}

func concat(store graphql.Fields, append graphql.Fields) graphql.Fields {
	for i, v := range append {
		store[i] = v
	}
	return store
}
