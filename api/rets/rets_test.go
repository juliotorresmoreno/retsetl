package rets

import (
	"testing"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"github.com/jpfielding/gorets/cmds/common"
)

func CompareResourceClassTest(t *testing.T) {
	t.Run("CompareResourceClassTest", func(t *testing.T) {
		mls := "5955ba4dc8cfd1262f107819"
		classID := "594ed02ec8cfd142614de0cd"
		resource := "Property"
		data := "../../data/flow.sqlite3"
		hostConn, err := db.OpenDSN(data)
		if err != nil {
			t.Error("1", err)
		}
		mlsConn, err := db.OpenMls(mls)
		if err != nil {
			t.Error("2", err)
		}
		err = CompareResourceClass(hostConn, mlsConn, mls, classID, resource)
		if err != nil {
			t.Error("3", err)
		}
	})
}

func TestMain(t *testing.T) {
	t.Run("GetResources", func(t *testing.T) {
		mls := "5955ba4dc8cfd1262f107819"
		if mlsConn, err := db.OpenMls(mls); err == nil {
			if err := GetResources(mlsConn, mls); err != nil {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	})
}

func GetResourcesTest(t *testing.T) {
	config := common.Config{
		Username:    "19147r",
		Password:    "pr7ph5fe",
		URL:         "http://imls.apps.retsiq.com/contact/rets/login",
		Version:     "RETS/1.0",
		UserAgent:   "",
		UserAgentPw: "",
	}
	meta := getMetadata(config)
	if meta.Rets.MetadataSystem.Version != config.Version {
		t.Errorf("getMetadata Error, version %v", meta.Rets.MetadataSystem.Version)
	}
}
