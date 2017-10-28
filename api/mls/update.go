package mls

import (
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/jpfielding/gorets/cmds/common"
	"strings"
)

//Update It updates in the database the parameters of configuration of the etl of the client
func Update(row *models.Mls) (bool, error) {
	_, err := helper.ValidateStruct(row)
	if err != nil && !strings.Contains(err.Error(), "password") {
		return false, err
	}
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		now := models.Mls{}
		_, err = conn.Id(row.ID).Get(&now)
		if err != nil {
			return false, err
		}

		if row.Password == "" {
			row.Password = now.Password
		}
		if row.PasswordBd == "" {
			row.PasswordBd = now.PasswordBd
		}
		row.CreateAt = now.CreateAt
		row.UpdateAt = now.UpdateAt

		if row.VersionRets == "Auto" {
			version, err := calculateMls(*row)
			if err != nil {
				return false, err
			}
			row.VersionRets = version
		} else {
			config := common.Config{
				Username:    row.Username,
				Password:    row.Password,
				URL:         row.URL,
				Version:     row.VersionRets,
				UserAgent:   row.UseragentName,
				UserAgentPw: row.UseragentPassword,
			}
			if ok, err := isValid(config); !ok {
				return false, err
			}
		}

		_, err = conn.ID(row.ID).Update(row)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err
	}
}
