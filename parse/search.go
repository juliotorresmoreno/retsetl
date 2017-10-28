package parse

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

//Search This object stores the list of results obtained from the search.
type Search struct {
	data []SearchElement
}

//SearchElement Represents a record of the search result
type SearchElement map[string]string

//LoadData You can parse the result and get the listing through the Search object
func (el *Search) LoadData(document *os.File) error {
	r := csv.NewReader(bufio.NewReader(document))
	record, err := r.Read()
	fields := record
	for err != io.EOF {
		record, err = r.Read()
		it := SearchElement{}
		limit := len(record)
		for i := 0; i < limit; i++ {
			it[fields[i]] = record[i]
		}
		el.data = append(el.data, it)
	}
	return nil
}
