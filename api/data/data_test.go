package data

import (
	"testing"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/go-xorm/xorm"
)

func TestMain(t *testing.T) {
	//t.Run("Import data", importDataTest)
	t.Run("Websockect", websockect)
}

func websockect(t *testing.T) {

}

func importDataTest(t *testing.T) {
	var hostConn, mlsConn, remoteConn *xorm.Engine
	var err error
	var mls, classID = "594566f3c8cfd11160b221ad", "595c1333c8cfd1506ef5c470"
	if hostConn, err = db.Open(); err != nil {
		t.Error(err)
		return
	}
	defer hostConn.Close()
	if mlsConn, err = db.OpenMls(mls); err != nil {
		t.Error(err)
		return
	}
	defer mlsConn.Close()

	row := models.Mls{}
	mlsConn.
		Where("id = ?", mls).
		Find(row)
	if remoteConn, err = db.OpenMlsConn(hostConn, mls); err != nil {
		t.Error(err)
		return
	}
	defer mlsConn.Close()
	if err = ImportData(hostConn, mlsConn, remoteConn, classID); err != nil {
		t.Error(err)
		return
	}
}
