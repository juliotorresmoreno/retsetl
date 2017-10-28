package db

import (
	"fmt"
	"strings"

	"bitbucket.org/mlsdatatools/retsetl/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

//OpenDSN Get the connection to the DB
func OpenDSN(dsn string) (*xorm.Engine, error) {
	orm, err := xorm.NewEngine(config.DRIVER, dsn)
	return orm, err
}

//Open Get the connection to the DB
func Open() (*xorm.Engine, error) {
	switch config.DRIVER {
	case "sqlite3":
		path := config.PATH
		return openSqlite(path + "/data/flow.sqlite3")
	case "mysql":
		fallthrough
	case "postgres":
		usr := config.DB_USER
		pwd := config.DB_PWD
		host := config.DB_HOST
		port := config.DB_PORT
		bd := config.DB_NAME
		orm, err := OpenBD(usr, pwd, host, port, bd, "postgres")
		if err != nil {
			fmt.Println(fmt.Errorf("connection error: %v", err))
		}
		return orm, err
	}
	return &xorm.Engine{}, fmt.Errorf("Driver not found")
}

func openSqlite(dsn string) (*xorm.Engine, error) {
	orm, err := xorm.NewEngine(config.DRIVER, dsn)
	if err != nil {
		fmt.Println(fmt.Errorf("connection error: %v", err))
	}
	return orm, err
}

//OpenMls Get the connection to the DB
func OpenMls(mls string) (*xorm.Engine, error) {
	path := fmt.Sprintf("%v/data/servers/%v.sqlite3", config.PATH, mls)
	return openSqlite(path)
}

//OpenMlsConn Get the connection to the DB
func OpenMlsConn(conn *xorm.Engine, mls string) (*xorm.Engine, error) {
	sql := fmt.Sprintf("select * from mls where id = '%v'", mls)
	result, err := conn.Query(sql)
	if err != nil {
		return &xorm.Engine{}, err
	}
	if len(result) == 0 {
		return &xorm.Engine{}, fmt.Errorf("Could not establish communication with the database")
	}
	usr := string(result[0]["username_bd"])
	pwd := string(result[0]["password_bd"])
	host := string(result[0]["server_bd"])
	bd := string(result[0]["name_bd"])
	port := ""
	if strings.Contains(host, ":") {
		tmp := strings.Split(host, ":")
		host = tmp[0]
		port = tmp[1]
	} else {
		port = "5432"
	}
	orm, err := OpenBD(usr, pwd, host, port, bd, "postgres")
	if err == nil {
		err = orm.Ping()
	}
	return orm, err
}

//OpenBD Get the connection to the DB
func OpenBD(usr, pwd, host, port, bd, driver string) (*xorm.Engine, error) {
	switch driver {
	case "mysql":
		charset := "?charset=utf8&parseTime=true"
		dsn := ""
		if config.DB_PWD != "" {
			dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v%v", usr, pwd, host, port, bd, charset)
		} else {
			dsn = fmt.Sprintf("%v@tcp(%v:%v)/%v%v", usr, host, host, port, charset)
		}
		return xorm.NewEngine(driver, dsn)
	case "postgres":
		dsn := ""
		if bd == "" {
			bd = "postgres"
		}
		if config.DB_PWD != "" {
			dsn = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", usr, pwd, host, port, bd)
		} else {
			dsn = fmt.Sprintf("user=%v dbname=%v host=%v port=%v password=paramore", usr, bd, host, port)
		}
		fmt.Println(dsn)
		return xorm.NewEngine(driver, dsn)
	}
	return nil, fmt.Errorf("Could not connect")
}
