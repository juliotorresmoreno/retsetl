package controllers

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/mlsdatatools/retsetl/api/mls"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
)

//GetMls Query a log of connection to a server rets
func GetMls(w http.ResponseWriter, r *http.Request) {
	result := make([]models.Mls, 0)
	mls.Search(&result)
	helper.Cors(w, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getMlsSuccess{
		Success: true,
		Data:    result,
	})
}

//PostMls Creating a log of connection to a server rets
func PostMls(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(int64(1000000))
	row := models.Mls{}
	row.Name = r.Form.Get("name")
	row.SystemID = r.Form.Get("system_id")
	row.URL = r.Form.Get("url")
	row.Username = r.Form.Get("username")
	row.Password = r.Form.Get("password")
	row.UseragentName = r.Form.Get("useragent_name")
	row.UseragentPassword = r.Form.Get("useragent_password")
	row.VersionRets = r.Form.Get("version")
	_, err := mls.Create(&row)
	helper.Cors(w, r)
	if err != nil {
		errs := err.(mls.Error)
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(postMlsError{
			Success: false,
			Errors:  errs.Details(),
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(postMlsSuccess{
		Success: true,
		Record:  row,
	})
}

//PutMls Updating a log of connection to a server rets
func PutMls(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(int64(1000 * 1000))
	row := models.Mls{}
	row.Name = r.Form.Get("name")
	row.SystemID = r.Form.Get("system_id")
	row.URL = r.Form.Get("url")
	row.Username = r.Form.Get("username")
	row.Password = r.Form.Get("password")
	row.UseragentName = r.Form.Get("useragent_name")
	row.UseragentPassword = r.Form.Get("useragent_password")
	row.VersionRets = r.Form.Get("version")
	_, err := mls.Create(&row)
	helper.Cors(w, r)
	if err != nil {
		errs := err.(mls.Error)
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(postMlsError{
			Success: false,
			Errors:  errs.Details(),
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(postMlsSuccess{
		Success: true,
		Record:  row,
	})
}

type getMlsSuccess struct {
	Success bool         `json:"success"`
	Data    []models.Mls `json:"data"`
}

type putMlsSuccess struct {
	Success bool       `json:"success"`
	Record  models.Mls `record:"record"`
}

type postMlsSuccess struct {
	Success bool       `json:"success"`
	Record  models.Mls `record:"record"`
}

type postMlsError struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Errors  models.Mls `json:"errors"`
}
