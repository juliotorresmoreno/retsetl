package mls

import (
	"time"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"
)

//Delete Delete in the database the client etl
func Delete(id string, row *models.Mls) (bool, error) {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		_, err := conn.Id(id).Get(row)
		if err != nil {
			return false, err
		}
		row.DeleteAt = time.Now()

		_, err = conn.ID(row.ID).Update(row)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err

	}
}
