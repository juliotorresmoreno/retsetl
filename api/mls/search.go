package mls

import (
	"time"

	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"
)

//Search Gets the configuration of the specified etl from its id or all if it is not specified.
func Search(row *[]models.Mls, id ...string) (bool, error) {
	if conn, err := db.Open(); err == nil {
		defer conn.Close()
		session := conn.NewSession()
		zero := time.Date(time.Now().Year()-100, 1, 1, 0, 0, 0, 0, time.UTC)
		session = session.Where("delete_at <= ?", zero)
		for _, v := range id {
			session = session.Where("id = ?", v)
		}
		err = session.Find(row)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return true, err
	}
}
