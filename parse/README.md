# README #

### What is this parser?

It is a tool that allows interpreting the results of the search and retrieval of data from the protocol rets

### Quickstart

For installation we proceed to execute:

go get -v bitbucket.org/mlsdatatools/retsetl/parse

### How is it used?

import "bitbucket.org/mlsdatatools/retsetl/parse"

#### To get the metadata we use

meta := "data/metadata.xml"
parser := NewParser()
var resp Meta
var err error
if resp, err = parser.OpenMeta(meta); err != nil {
    t.Error(err)
}
//resp Contains an object with the result that can be queried

#### To get the search we use

search := "data/search.csv"
parser := NewParser()
var resp Search
var err error
if resp, err = parser.OpenSearch(search); err != nil {
    t.Error(err)
}


### Test results

To perform the test of the application we must use

go test -v

![Result](result.png?raw=true "Result")