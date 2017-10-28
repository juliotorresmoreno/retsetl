package taskmgr

import (
	"fmt"
	"strings"
	"time"

	"log"

	"bitbucket.org/mlsdatatools/retsetl/api/data"
	"bitbucket.org/mlsdatatools/retsetl/db"
	"bitbucket.org/mlsdatatools/retsetl/models"
)

var tasks chan task

type task struct {
	mls      string
	class    string
	schedule string
}

func Start() {
	tasks = make(chan task)
	go exec()
	go finder()
}
func exec() {
	var _task task
	hostConn, err := db.Open()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer hostConn.Close()
	for {
		_task = <-tasks
		func() {
			if mlsConn, err := db.OpenMls(_task.mls); err == nil {
				remoteConn, err := db.OpenMlsConn(hostConn, _task.mls)
				defer mlsConn.Close()
				if err == nil {
					row := models.Schedule{
						Status: "Importing",
					}
					_, err = hostConn.ID(_task.schedule).Cols("status").Update(row)
					if err == nil {
						data.ImportData(hostConn, mlsConn, remoteConn, "", _task.class)
					}
				}
			}
		}()
	}
}

func finder() {
	hostConn, err := db.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer hostConn.Close()
	for {
		now := time.Now()
		day := strings.ToLower(fmt.Sprintf("%v", now.Weekday()))
		hour := fmt.Sprintf("%v", now.Hour())
		if len(hour) == 1 {
			hour = "0" + hour
		}
		min := fmt.Sprintf("%v", now.Minute())
		if len(min) == 1 {
			min = "0" + min
		}
		_time := fmt.Sprintf("%v:%v:00", hour, min)
		sql := fmt.Sprintf(
			`SELECT * 
               FROM schedule 
			  WHERE classes != '' 
				AND classes IS NOT NULL
				AND hour = '%v'
				AND (frequency = 1 or frequency = 2 or (%v = 1 and frequency = 3))`,
			_time, day,
		)
		if result, err := hostConn.QueryString(sql); err == nil {
			for _, v := range result {
				classes := strings.Split(v["classes"], ",")
				for _, class := range classes {
					tasks <- task{
						schedule: v["id"],
						mls:      v["mls"],
						class:    class,
					}
				}
			}
		} else {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Minute)
	}
}
