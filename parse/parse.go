package parse

import (
	"io/ioutil"
	"os"
)

//Parser structure responsible for performing the parse of the document
type Parser struct {
	document string
}

//NewParser Constructor of the new parseo object
func NewParser() Parser {
	return Parser{}
}

//OpenMeta Open a metadata response document
func (el Parser) OpenMeta(file string) (Meta, error) {
	var err error
	var document []byte
	var response Meta
	document, err = ioutil.ReadFile(file)
	if err != nil {
		return response, err
	}
	response.LoadData(document)
	return response, err
}

//OpenSearch Open a search response document
func (el Parser) OpenSearch(file string) (Search, error) {
	var err error
	var response Search
	var rfile *os.File
	rfile, err = os.Open(file)
	if err != nil {
		return response, err
	}
	response.LoadData(rfile)
	return response, err
}
