package helper

import (
	"github.com/graphql-go/graphql"
)

//GetGraphParam Get the parameter of a graphql request securely
//@param params: Contains all the information coming from the user
//@param name:   Name of the parameter to be extracted
//@param zero:   Default value
func GetGraphParam(params graphql.ResolveParams, name string, zero interface{}) interface{} {
	result, ok := params.Args[name]
	if ok && result != nil {
		return result
	}
	return zero
}
