package parse

import (
	"fmt"
	"strings"

	"github.com/clbanning/mxj"
)

//Meta represents the metadata obtained from the response
type Meta struct {
	Rets struct {
		MetadataSystem struct {
			SystemID          string
			SystemDescription string
			Version           string
			Date              string
		}
		MetadataResource    []MetadataElement
		MetadataClass       []MetadataElement
		MetadataTable       []MetadataElement
		MetadataLookupType  []MetadataElement
		MetadataForeignKeys []MetadataElement
	}
	document mxj.Map
}

//MetadataElement Represents a record of the metadata itself
type MetadataElement struct {
	Class         string
	Resource      string
	Date          string
	Version       string
	columnsID     map[string]int
	Columns       []string
	columnsLength int
	Data          [][]string
	dataLength    int
}

func (el MetadataElement) Get(row int, key string) string {
	value, ok := el.columnsID[key]
	if !ok || value >= len(el.Data[row]) {
		return ""
	}
	return el.Data[row][value]
}

//LoadData Generates structure data from a given document
func (el *Meta) LoadData(document []byte) error {
	var result mxj.Map
	var err error
	result, err = mxj.NewMapXml(document)
	if err != nil {
		return err
	}

	el.document = result

	el.Rets.MetadataSystem.Date = getValueStr(result, "RETS.METADATA-SYSTEM.-Date")
	el.Rets.MetadataSystem.Version = getValueStr(result, "RETS.METADATA-SYSTEM.-Version")
	el.Rets.MetadataSystem.SystemID = getValueStr(result, "RETS.METADATA-SYSTEM.SYSTEM.-SystemID")
	el.Rets.MetadataSystem.SystemDescription = getValueStr(result, "RETS.METADATA-SYSTEM.SYSTEM.-SystemDescription")

	el.Rets.MetadataResource = el.getMetadata("RESOURCE")
	el.Rets.MetadataClass = el.getMetadata("CLASS")
	el.Rets.MetadataTable = el.getMetadata("TABLE")
	el.Rets.MetadataLookupType = el.getMetadata("LOOKUP_TYPE")
	el.Rets.MetadataForeignKeys = el.getMetadata("FOREIGNKEYS")
	return nil
}

func (el *Meta) getMetadata(resource string) []MetadataElement {
	c, _ := el.document.ValuesForPath(fmt.Sprintf("RETS.METADATA-%v", resource))

	r := make([]MetadataElement, 0)

	for _, t := range c {
		v := t.(map[string]interface{})
		columnsStr := v["COLUMNS"].(string)
		columnsSpl := strings.Split(columnsStr, "\t")
		Elm := MetadataElement{
			Columns:  columnsSpl,
			Class:    getStringValue(v["-Class"]),
			Resource: getStringValue(v["-Resource"]),
			Date:     getStringValue(v["-Date"]),
			Version:  getStringValue(v["-Version"]),
		}
		Elm.columnsLength = len(columnsSpl)
		Elm.columnsID = map[string]int{}
		for k, v := range columnsSpl {
			Elm.columnsID[v] = k
		}
		switch fmt.Sprintf("%T", v["DATA"]) {
		case "[]interface {}":
			dataElemts := v["DATA"].([]interface{})
			for _, dataStr := range dataElemts {
				dataSpl := strings.Split(dataStr.(string), "\t")
				Elm.Data = append(Elm.Data, dataSpl)
				Elm.dataLength = len(dataSpl)
			}
		case "string":
			dataElemts := v["DATA"].(string)
			dataSpl := strings.Split(dataElemts, "\t")
			Elm.Data = append(Elm.Data, dataSpl)
			Elm.dataLength = len(dataSpl)
		}
		r = append(r, Elm)
	}
	return r
}

func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}

func getValueStr(xml mxj.Map, name string) string {
	val, _ := xml.ValueForPath(name)
	return val.(string)
}
