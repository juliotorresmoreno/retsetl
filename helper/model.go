package helper

import (
	"regexp"

	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
)

func init() {
	alphaspacesnum, _ := regexp.Compile("^[a-zA-Z][a-zA-Z0-9( )]*$")
	govalidator.TagMap["alphaspacesnum"] = govalidator.Validator(func(str string) bool {
		return alphaspacesnum.MatchString(str)
	})
	govalidator.TagMap["rets"] = govalidator.Validator(func(str string) bool {
		versions := []string{"RETS/1.0", "RETS/1.5", "RETS/1.7.2", "RETS/1.8.1", "Auto"}
		return govalidator.IsIn(str, versions...)
	})
	govalidator.TagMap["server"] = govalidator.Validator(func(str string) bool {
		return true
	})
}

//ValidateStruct Performs the validations of the fields in a model of na table in database
func ValidateStruct(s interface{}) (bool, error) {
	return govalidator.ValidateStruct(s)
}

//Sync Creates the necessary tables to run the application in an
//external database created to store the metadata
func Sync(conn *xorm.Engine) {
	conn.Sync2(models.Class{})
	conn.Sync2(models.Resource{})
	conn.Sync2(models.Table{})
}
