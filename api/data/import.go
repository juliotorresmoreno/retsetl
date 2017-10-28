package data

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/mlsdatatools/retsetl/config"

	"gopkg.in/mgo.v2/bson"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"

	"strings"

	"encoding/json"

	"strconv"

	"bitbucket.org/mlsdatatools/retsetl/api/mail"
	"bitbucket.org/mlsdatatools/retsetl/api/metadata"
	"bitbucket.org/mlsdatatools/retsetl/api/mls"
	configuration "bitbucket.org/mlsdatatools/retsetl/config"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"bitbucket.org/mlsdatatools/retsetl/ws"
	"github.com/go-xorm/xorm"
)

func mapColumns(field map[string]interface{}, columns map[string]string) (string, string, bool) {
	_c := []string{
		string(field["long_name"].([]uint8)),
		string(field["short_name"].([]uint8)),
		string(field["system_name"].([]uint8)),
		string(field["standard_name"].([]uint8)),
	}
	for _, v := range _c {
		if t, ok := columns[v]; v != "" && ok {
			return _c[2], t, true
		}
	}
	return "", "", false
}

func mapField(columns []map[string]interface{}) map[string]string {
	result := map[string]string{}
	for _, v := range columns {
		c := strings.Split("Synonym", ",")
		s := string(v["StandardName"].([]uint8))
		for _, x := range c {
			result[x] = s
		}
		result[s] = s
	}
	return result
}

//ImportData Extracts data from the rest server and stores them in bd
func ImportData(hostConn, mlsConn, remoteConn *xorm.Engine, clientID, classID string) error {
	class := models.Class{}
	cols := make([]map[string]interface{}, 0)
	mlsConn.Id(classID).Get(&class)
	if err := mlsConn.Select("system_name, standard_name, long_name, short_name").
		Table(models.Table{}.TableName()).
		Where("class_id = ?", classID).
		Find(&cols, models.Table{ClassID: classID}); err != nil {
		return err
	}

	importLog := models.ImportLog{
		ID:         bson.NewObjectId().Hex(),
		Class:      class.StandardName,
		ClassID:    class.ID,
		Resource:   class.ResourceName,
		ResourceID: class.ResourceID,
		Mls:        class.Mls,
		Date:       time.Now(),
		Status:     "PROCCESS",
	}
	_, err := hostConn.Insert(&importLog)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}

	fields := make([]map[string]interface{}, 0)
	hostConn.Select("StandardName, Synonym").
		Table(class.StoreAs).
		Where("Synonym <> ''").
		Find(&fields)
	_mapFields := mapField(fields)

	_columns := map[string]string{}
	columnsNames := make([]columns, 0)
	for _, v := range cols {
		if t, s, ok := mapColumns(v, _mapFields); ok {
			columnsNames = append(columnsNames, columns{SystemName: s})
			_columns[t] = s
		}
	}
	_columns["L_UpdateDate"] = "UpdateDate"
	tablename := class.StoreAs
	remoteSession := orm.NewSession(remoteConn)
	defer remoteSession.Close()
	row := orm.Model{}
	row.SetTableName(tablename)
	createContext(remoteSession, row, columnsNames)

	sql := fmt.Sprintf(
		`SELECT max("UpdateDate") AS last 
		   FROM "%v"
		  WHERE "Source" = '%v'`,
		class.ResourceName, class.ID,
	)
	last, err := remoteConn.QueryString(sql)
	if err != nil {
		return err
	}
	_last := ""
	if _, ok := last[0]["last"]; ok {
		if _last = last[0]["last"]; _last != "" {
			year, _ := strconv.Atoi(_last[0:4])
			month, _ := strconv.Atoi(_last[5:7])
			day, _ := strconv.Atoi(_last[8:10])
			hour, _ := strconv.Atoi(_last[11:13])
			minute, _ := strconv.Atoi(_last[14:16])
			second, _ := strconv.Atoi(_last[17:19])
			t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
			t = t.Add(1 * time.Second)
			_last = fmt.Sprintf("%v-%v-%vT%v:%v:%v", year, month, day, hour, minute, second)
		}
	}

	imported, err := searchRets(
		hostConn,
		mlsConn,
		remoteConn,
		clientID,
		classID,
		tablename,
		importLog.ID,
		_columns,
		_last,
	)
	importLog.Status = "COMPLETE"
	importLog.Imported = fmt.Sprintf("%v", imported)
	hostConn.Id(importLog.ID).Cols("status", "imported").Update(importLog)
	sendMail(class, importLog, err)
	return err
}

func sendMail(class models.Class, importLog models.ImportLog, err error) {
	var subject = "Import Report"
	var to = config.EMAIL_ADMIN
	var from = config.EMAIL_SEND
	var pass = config.EMAIL_PASSWORD
	var body = ""
	var className = fmt.Sprintf("%v (%v)", class.StandardName, class.ClassName)
	var resource = class.ResourceName
	var imported = importLog.Imported
	var url = fmt.Sprintf(
		"http://199.168.136.144:8080/logs/import/%v/%v",
		importLog.Mls, importLog.ID,
	)
	if err != nil {
		str := "An error occurred while importing class %v from resource %v\n\n"
		str += "We were able to import %v records\n\n"
		str += "To see more details, access	%v\n\n"
		str += "Error details: %v"
		body = fmt.Sprintf(str, className, resource, imported, url, err)
	} else {
		str := "Has successfully imported class %v from resource%v\n\n"
		str += "We were able to import %v records\n\n"
		str += "To see more details, access	%v\n\n"
		body = fmt.Sprintf(str, className, resource, imported, url)
	}
	mail.Send(subject, to, from, pass, body)
}

func searchRets(
	hostConn, mlsConn, remoteConn *xorm.Engine,
	clientID, classID, tablename, importLog string,
	columns map[string]string, _last string,
) (int, error) {
	rows := make([]models.Mls, 0)
	meta := models.Class{}
	resource := models.Resource{}
	metadata.GetClassByID(mlsConn, &meta, classID)
	metadata.GetResourcesByID(mlsConn, &resource, meta.ResourceID)
	mls.Search(&rows, resource.Mls)
	if len(rows) == 0 {
		importLogDetails := models.ImportLogDetails{
			ID:          bson.NewObjectId().Hex(),
			ImportLogID: importLog,
			Data:        "",
			Error:       "Not found mls",
		}
		hostConn.Insert(importLogDetails)
		return 0, fmt.Errorf("Not found mls")
	}
	config := common.Config{
		Username:    rows[0].Username,
		Password:    rows[0].Password,
		URL:         rows[0].URL,
		Version:     rows[0].VersionRets,
		UserAgent:   rows[0].UseragentName,
		UserAgentPw: rows[0].UseragentPassword,
	}
	if config.UserAgent == "" {
		config.UserAgent = "Threewide/1.0"
	}

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		importLogDetails := models.ImportLogDetails{
			ID:          bson.NewObjectId().Hex(),
			ImportLogID: importLog,
			Data:        "",
			Error:       err.Error(),
		}
		hostConn.Insert(importLogDetails)
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	urls, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err != nil {
		return 0, err
	}
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: urls.Logout})

	if urls.GetPayloadList != "" {
		fmt.Println("Payloads: ", urls.GetPayloadList)
		payloads, err := rets.GetPayloadList(ctx, session, rets.PayloadListRequest{
			URL: urls.GetPayloadList,
			ID:  fmt.Sprintf("%s:%s", resource.ResourceID, meta.ClassName),
		})

		if err != nil {
			importLogDetails := models.ImportLogDetails{
				ID:          bson.NewObjectId().Hex(),
				ImportLogID: importLog,
				Data:        "",
				Error:       err.Error(),
			}
			hostConn.Insert(importLogDetails)
			return 0, err
		}
		err = payloads.ForEach(func(payload rets.CompactData, err error) error {
			fmt.Printf("%v\n", payload)
			importLogDetails := models.ImportLogDetails{
				ID:          bson.NewObjectId().Hex(),
				ImportLogID: importLog,
				Data:        "",
				Error:       err.Error(),
			}
			hostConn.Insert(importLogDetails)
			return err
		})
		if err != nil {
			importLogDetails := models.ImportLogDetails{
				ID:          bson.NewObjectId().Hex(),
				ImportLogID: importLog,
				Data:        "",
				Error:       err.Error(),
			}
			hostConn.Insert(importLogDetails)
			return 0, err
		}
	}

	fmt.Println("Search: ", urls.Search)
	query := "*"
	if _last != "" {
		query = fmt.Sprintf("(L_UpdateDate=%v+)", _last)
	}
	_columns := make([]string, 0)
	for k := range columns {
		_columns = append(_columns, k)
	}
	req := rets.SearchRequest{
		URL: urls.Search,
		SearchParams: rets.SearchParams{
			Select:     strings.Join(_columns, ","),
			Query:      query,
			SearchType: resource.ResourceID,
			Class:      meta.ClassName,
			Format:     "",
			QueryType:  "dmql2",
			Count:      0,
			Limit:      configuration.LIMIT,
			Offset:     0,
		},
	}
	imported := processCompact(
		ctx,
		hostConn,
		remoteConn,
		session,
		req,
		classID,
		clientID,
		tablename,
		importLog,
		columns,
	)
	return imported, nil
}

func processCompact(
	ctx context.Context,
	hostConn, remoteConn *xorm.Engine,
	sess rets.Requester, req rets.SearchRequest,
	classID, clientID, tablename, importLog string,
	columns map[string]string,
) int {
	// loop over all the pages we need
	session := orm.NewSession(remoteConn)
	hub := ws.GetHub()
	for {
		data := make([]orm.Model, 0)
		fmt.Printf("Querying next page\n")
		result, err := rets.SearchCompact(ctx, sess, req)
		if err != nil {
			importLogDetails := models.ImportLogDetails{
				ID:          bson.NewObjectId().Hex(),
				ImportLogID: importLog,
				Data:        "",
				Error:       err.Error(),
			}
			hostConn.Insert(importLogDetails)
		}
		switch result.Response.Code {
		case rets.StatusOK:
		case rets.StatusNoRecords:
			return req.Offset
		case rets.StatusSearchError:
			fallthrough
		default: // shit hit the fan
			importLogDetails := models.ImportLogDetails{
				ID:          bson.NewObjectId().Hex(),
				ImportLogID: importLog,
				Data:        "",
				Error:       result.Response.Text,
			}
			hostConn.Insert(importLogDetails)
		}
		count := 0
		hasMoreRows, err := result.ForEach(func(row rets.Row, err error) error {
			if err != nil {
				importLogDetails := models.ImportLogDetails{
					ID:          bson.NewObjectId().Hex(),
					ImportLogID: importLog,
					Data:        "",
					Error:       err.Error(),
				}
				hostConn.Insert(importLogDetails)
				fmt.Printf("Ups, %v", err)
				return err
			}
			record := orm.Model{}
			record.SetTableName(tablename)
			for index, col := range result.Columns {
				if index < len(row) {
					if _col, ok := columns[col]; ok {
						record[_col] = row[index]
					}
				}
			}
			record["Source"] = classID
			data = append(data, record)
			count++
			return err
		})
		result.Close()
		if err != nil {
			importLogDetails := models.ImportLogDetails{
				ID:          bson.NewObjectId().Hex(),
				ImportLogID: importLog,
				Data:        "",
				Error:       err.Error(),
			}
			hostConn.Insert(importLogDetails)
		}
		if !hasMoreRows {
			return req.Offset
		}
		if req.Offset == 0 {
			req.Offset = 1
		}

		session.InsertAll(data)
		req.Offset = req.Offset + count
		message, _ := json.Marshal(map[string]interface{}{
			"type":     "websocket/PROGRESS",
			"progress": req.Offset,
		})
		hub.Send(clientID, message)
	}
}

type columns struct {
	SystemName string
}

func (el columns) TableName() string {
	return models.Table{}.TableName()
}
