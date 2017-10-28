package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"

	"github.com/kardianos/osext"
)

//PORT Http server listening port
var PORT string

//DB_HOST Mongo database server
var DB_HOST string

//DB_PORT Port of the Mongo database
var DB_PORT string

//DB_NAME Name of the mongo database
var DB_NAME string

//DB_USER Port of the Mongo database
var DB_USER string

//DB_PWD Name of the mongo database
var DB_PWD string

//DRIVER Name of the mongo database
var DRIVER string

//READ_TIMEOUT Tiempo de espera que tardara el servidor en recibir los datos
var READ_TIMEOUT time.Duration

//LIMIT d
var LIMIT int

//EMAIL_ADMIN d
var EMAIL_ADMIN string

//EMAIL_SEND d
var EMAIL_SEND string

//EMAIL_PASSWORD d
var EMAIL_PASSWORD string

//PATH d
var PATH string

func init() {
	executableFolder, _ := osext.ExecutableFolder()
	path := flag.String("path", executableFolder, "")
	flag.Parse()
	PATH = *path
	text, err := ioutil.ReadFile(PATH + "/config/config.json")
	if err != nil {
		panic(err)
	}
	var data = &configuration{}
	err = json.Unmarshal(text, data)
	if err != nil {
		panic(err)
	}
	PORT = data.Port
	READ_TIMEOUT = data.ReadTimeout

	DB_HOST = data.DbHost
	DB_PORT = data.DbPort
	DB_NAME = data.DbName
	DB_USER = data.DbUser
	DB_PWD = data.DbPwd
	DRIVER = data.Driver

	EMAIL_ADMIN = data.EmailAdmin
	EMAIL_SEND = data.EmailSend
	EMAIL_PASSWORD = data.EmailPassword

	LIMIT = data.Limit
}
