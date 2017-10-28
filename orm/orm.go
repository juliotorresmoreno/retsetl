package orm

import (
	"fmt"

	"bitbucket.org/mlsdatatools/retsetl/config"
)

type datatypes map[string]datatype

type datatype struct {
	MySQL    string
	Postgres string
	Oracle   string
	SQLSvr   string
	SQLite3  string
}

var relacion = datatypes{
	"text": datatype{
		MySQL:    "text",
		Postgres: "text",
		Oracle:   "",
		SQLSvr:   "",
		SQLite3:  "text",
	},
}

var quotes = map[string]string{
	"mysql":    "'",
	"postgres": "'",
	"oracle":   "",
	"sqlsvr":   "",
}

var separator = map[string]string{
	"mysql":    "`",
	"postgres": "\"",
	"oracle":   "",
	"sqlsvr":   "",
}

func getType(str string) (string, error) {
	if val, ok := relacion[str]; ok {
		switch config.DRIVER {
		case "mysql":
			return val.MySQL, nil
		case "postgres":
			return val.Postgres, nil
		case "sqlite3":
			return val.SQLite3, nil
		case "oracle":
			return val.Oracle, nil
		case "sqlsvr":
			return val.SQLSvr, nil
		}
	}
	return "", fmt.Errorf("Field datatype not found")
}
