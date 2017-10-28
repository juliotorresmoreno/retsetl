package mls

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
	"gopkg.in/mgo.v2/bson"
)

//Error Custom error that allows you to obtain additional or more detailed information.
type Error struct {
	err     string
	details models.Mls
}

func (el Error) Error() string {
	return el.err
}

//Details Returns additional and detailed information of the errors occurred
func (el Error) Details() models.Mls {
	return el.details
}

//Create Insert in the database the new client etl
func Create(row *models.Mls) (bool, error) {
	_, err := helper.ValidateStruct(row)
	if err != nil {
		return false, err
	}
	if hostConn, err := db.Open(); err == nil {
		if row.ID == "" {
			row.ID = bson.NewObjectId().Hex()
		}
		row.DeleteAt = time.Date(time.Now().Year()-100, 1, 1, 0, 0, 0, 0, time.UTC)
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
		_, err = hostConn.InsertOne(row)
		mlsCon, _ := db.OpenMls(row.ID)
		helper.Sync(mlsCon)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err
	}
}

func calculateMls(mls models.Mls) (string, error) {
	versions := []string{"RETS/1.0", "RETS/1.5", "RETS/1.7.2"}
	for _, version := range versions {
		config := common.Config{
			Username:    mls.Username,
			Password:    mls.Password,
			URL:         mls.URL,
			Version:     version,
			UserAgent:   mls.UseragentName,
			UserAgentPw: mls.UseragentPassword,
		}
		if ok, _ := isValid(config); ok {
			return version, nil
		}
	}
	return "", fmt.Errorf("No valid server RETS")
}

func isValid(config common.Config) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if config.UserAgent == "" {
		config.UserAgent = "Threewide/1.0"
	}
	session, err := config.Initialize()
	if err != nil {
		return false, err
	}
	urls, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err == nil {
		defer rets.Logout(ctx, session, rets.LogoutRequest{URL: urls.Logout})
		return true, nil
	}
	return false, err
}
