package rets

import (
	"encoding/json"
	"fmt"
	"testing"

	"bitbucket.org/mlsdatatools/retsetl/parse"
	"github.com/jpfielding/gorets/cmds/common"
)

//TestGetMetadata Test of the GetMetadata function of the parser object
func TestGetMetadata(t *testing.T) {
	config := common.Config{
		Username: "19147r",
		Password: "pr7ph5fe",
		URL:      "http://imls.apps.retsiq.com/contact/rets/login",
		Version:  "RETS/1.0",
	}
	document, _ := GetMetadata(config)
	meta := parse.Meta{}
	meta.LoadData(document)
	rets, _ := json.Marshal(meta.Rets)
	fmt.Println(string(rets))
}
