package schedule

import (
	"bitbucket.org/mlsdatatools/retsetl/models"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2/bson"
)

//CreateSchedule Creates a schedule assignment in the system
func CreateSchedule(hostConn *xorm.Engine, mls string,
	sunday, monday, tuesday, wednesday, thursday, friday, saturday bool,
	frequency int,
	hour, classes string) error {
	row := models.Schedule{
		ID:        bson.NewObjectId().Hex(),
		Monday:    boolToInt(monday),
		Tuesday:   boolToInt(tuesday),
		Wednesday: boolToInt(wednesday),
		Thursday:  boolToInt(thursday),
		Friday:    boolToInt(friday),
		Saturday:  boolToInt(saturday),
		Sunday:    boolToInt(sunday),
		Frequency: frequency,
		Hour:      hour,
		Classes:   classes,
		Mls:       mls,
		Status:    "Pending",
	}
	if _, err := hostConn.Insert(row); err != nil {
		return err
	}
	return nil
}

//UpdateSchedule Update a schedule assignment in the system
func UpdateSchedule(hostConn *xorm.Engine,
	ID, mls string,
	sunday, monday, tuesday, wednesday, thursday, friday, saturday bool,
	frequency int,
	hour, classes string) error {
	row := models.Schedule{
		ID:        ID,
		Monday:    boolToInt(monday),
		Tuesday:   boolToInt(tuesday),
		Wednesday: boolToInt(wednesday),
		Thursday:  boolToInt(thursday),
		Friday:    boolToInt(friday),
		Saturday:  boolToInt(saturday),
		Sunday:    boolToInt(sunday),
		Frequency: frequency,
		Hour:      hour,
		Classes:   classes,
		Mls:       mls,
	}
	if _, err := hostConn.ID(ID).Cols(
		"monday", "tuesday", "wednesday", "thursday", "friday",
		"saturday", "sunday", "frequency", "hour", "classes", "mls",
	).Update(row); err != nil {
		return err
	}
	return nil
}

func boolToInt(data bool) int {
	if data {
		return 1
	}
	return 0
}

//DeleteSchedule Delete a schedule assignment in the system
func DeleteSchedule(hostConn *xorm.Engine, ID string) error {
	row := models.Schedule{ID: ID}
	if _, err := hostConn.Delete(row); err != nil {
		return err
	}
	return nil
}
