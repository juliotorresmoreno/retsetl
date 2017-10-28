package mls

import (
	"encoding/json"
	"fmt"
	"testing"

	"bitbucket.org/mlsdatatools/retsetl/models"
)

func TestUpdatePassword(t *testing.T) {
	mls := "59404e82c8cfd14190ccdec3"
	row := models.Mls{
		ID:       mls,
		Password: "pr7ph5fe",
	}
	Update(&row)
	rows := make([]models.Mls, 0)
	Search(&rows, mls)
	if len(rows) == 0 {
		t.Error("record not fount")
		return
	}
	if rows[0].Password != row.Password {
		t.Error("record not update")
	}
	data, _ := json.Marshal(rows[0])
	fmt.Println(string(data))
}
