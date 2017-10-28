package dictionaryGraphql

import (
	"testing"

	"bitbucket.org/mlsdatatools/retsetl/db"
)

func TestMain(t *testing.T) {
	conn, _ := db.Open()
	cleanFields(conn, "Property", "City")
}
