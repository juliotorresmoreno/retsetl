package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

//TestOpenMeta Test of the OpenMeta function of the parser object
func TestOpenMeta(t *testing.T) {
	meta := "data/metadata.xml"
	parser := NewParser()
	var resp Meta
	var err error
	if resp, err = parser.OpenMeta(meta); err != nil {
		t.Error(err)
	}
	if v := "IMLS"; resp.Rets.MetadataSystem.SystemID != v {
		t.Error(fmt.Errorf("Unexpected result, expected %v, found %v", v, resp.Rets.MetadataSystem.SystemID))
	}
	if v := "Intermountain MLS, Inc."; resp.Rets.MetadataSystem.SystemDescription != v {
		t.Error(fmt.Errorf("Unexpected result, expected %v, found %v", v, resp.Rets.MetadataSystem.SystemDescription))
	}
	data, _ := json.Marshal(resp)
	ioutil.WriteFile("data/metadata.json", data, 0777)
}

//TestOpenSearch Test of the OpenSearch function of the parser object
func TestOpenSearch(t *testing.T) {
	search := "data/search.csv"
	parser := NewParser()
	var resp Search
	var err error
	if resp, err = parser.OpenSearch(search); err != nil {
		t.Error(err)
	}
	data, _ := json.Marshal(resp.data)
	ioutil.WriteFile("data/search.json", data, 0777)
}
