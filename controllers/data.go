package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"github.com/go-xorm/builder"
)

//GetData Query any table with any number of columns
//This method can be used to get the data of any table
func GetData(w http.ResponseWriter, r *http.Request) {
	var err error
	var data []map[string]interface{}
	var values = r.URL.Query()
	var model = values["model"]
	var mls = values["mls"]
	helper.Cors(w, r)
	if len(mls) != 1 {
		err = fmt.Errorf("mls not found")
		responseError(w, err)
		return
	}
	if len(model) != 1 {
		err = fmt.Errorf("model not found")
		responseError(w, err)
		return
	}

	row := orm.Model{}
	row.SetTableName(model[0])
	cond := builder.Eq{}

	for k, v := range values {
		if k != "model" && k != "mls" {
			cond[k] = v[0]
		}
	}

	conn, err := db.OpenMls(mls[0])
	if err != nil {
		responseError(w, err)
		return
	}
	session := orm.NewSession(conn)
	data, err = session.Find(row, cond)

	if err != nil {
		responseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func responseError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
