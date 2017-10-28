package rets

import (
	"fmt"

	"time"

	"strings"

	"bitbucket.org/mlsdatatools/retsetl/api/metadata"
	"bitbucket.org/mlsdatatools/retsetl/api/mls"
	"bitbucket.org/mlsdatatools/retsetl/helper"
	"bitbucket.org/mlsdatatools/retsetl/models"
	"bitbucket.org/mlsdatatools/retsetl/orm"
	"bitbucket.org/mlsdatatools/retsetl/parse"
	"bitbucket.org/mlsdatatools/retsetl/rets"
	"github.com/go-xorm/xorm"
	"github.com/jpfielding/gorets/cmds/common"
	"gopkg.in/mgo.v2/bson"
)

var resourcesCollection = models.Resource{}.TableName()

func getFieldsDic(hostConn *xorm.Engine, resourceName string) ([]map[string]string, error) {
	fieldsDic := make([]map[string]string, 0)
	if err := hostConn.
		Select("StandardName, SimpleDataType, SugMaxLength, Synonym").
		Table(resourceName).
		Find(&fieldsDic); err != nil {
		return nil, err
	}
	return fieldsDic, nil
}

func getFieldsMeta(mlsConn *xorm.Engine, classID string) ([]map[string]string, error) {
	fieldsMeta := make([]map[string]string, 0)
	table := models.Table{}.TableName()
	if err := mlsConn.
		Select(`
			id,
			standard_name,
			data_type,
			maximum_length,
			class_name,
			system_name,
			short_name,
			long_name,
			class_id,
			resource_id
		`).
		Table(table).
		Where("class_id = ?", classID).
		Find(&fieldsMeta); err != nil {
		return nil, err
	}
	return fieldsMeta, nil
}

func compareFields(meta map[string]string, dic map[string]string, fields []string, max int) (map[string]bool, int, bool) {
	cnt := max
	compare := map[string]bool{}
	found := false
	_value := dic["StandardName"]
	values := strings.Split(dic["Synonym"], ",")
	values = append(values, _value)
	for _, field := range fields {
		value := meta[field]
		if value != "" {
			for _, _value := range values {
				if strings.ToLower(value) == strings.ToLower(_value) {
					found = true
					compare[meta[field]] = true
					cnt--
					break
				} else {
					compare[meta[field]] = false
				}
			}
		} else {
			cnt--
		}
	}
	return compare, cnt, found
}

func registerLog(hostConn *xorm.Engine, mls, resourceName, classID string) (models.Log, error) {
	log := models.Log{}
	log.ID = bson.NewObjectId().Hex()
	log.Date = time.Now()
	log.Action = "transform"
	log.ResourceName = resourceName
	log.ClassID = classID
	log.Mls = mls
	_, err := hostConn.Insert(log)
	return log, err
}

func removeMapped(mlsSession *orm.Session, field string, meta map[string]string) error {
	table := models.Table{}.TableName()
	tableField := orm.Model{
		"mapped":        "",
		"mapped_status": "",
	}
	tableField.SetTableName(table)
	cond := fmt.Sprintf("mapped = '%v' and class_id = '%v'", field, meta["class_id"])
	return mlsSession.Update(tableField, cond)
}

func assignMapped(mlsSession *orm.Session, field, fieldID string, meta map[string]string, status string) error {
	table := models.Table{}.TableName()
	tableField := orm.Model{
		"mapped":        field,
		"mapped_status": status,
	}
	tableField.SetTableName(table)
	cond := fmt.Sprintf("id = '%v'", fieldID)
	return mlsSession.Update(tableField, cond)
}

//CompareResourceClass Makes a comparison between the fields that are in the metadata
//and the fields that are in the dictionary, from here you will know which ones match
//and which do not.
//The comparison is made on the type and length and between the standardName of the
//dictionary vs. the system_name, standard_name, long_name, short_name
func CompareResourceClass(hostConn, mlsConn *xorm.Engine, mls, classID, resourceName string) error {
	//Verify if there is connection to the database
	if err := hostConn.Ping(); err != nil {
		return fmt.Errorf("1000 %v", err)
	}
	//Get dictionary fields
	fieldsDic, err := getFieldsDic(hostConn, resourceName)
	if err != nil {
		return fmt.Errorf("1010 %v", err)
	}
	//Get metadata fields
	fieldsMeta, err := getFieldsMeta(mlsConn, classID)
	if err != nil {
		return fmt.Errorf("1020 %v", err)
	}

	fieldsValidate := []string{"standard_name", "short_name", "long_name"}
	fieldsValidateLength := len(fieldsValidate)

	maxInsert := 20
	details := make([]models.LogDetails, 0, maxInsert)
	log, err := registerLog(hostConn, mls, resourceName, classID)
	if err != nil {
		return fmt.Errorf("1030 %v", err)
	}
	hostSession := orm.NewSession(hostConn)
	mlsSession := orm.NewSession(mlsConn)

	for _, fieldMeta := range fieldsMeta {
		_found := false
		for _, fieldDic := range fieldsDic {
			field := fieldDic["StandardName"]
			compare, cnt, found := compareFields(fieldMeta, fieldDic, fieldsValidate, fieldsValidateLength)
			if found {
				_found = true
				fieldDicMaxLength := fieldDic["StandardName"]
				fieldMetaMaxLength := fieldMeta["maximum_length"]
				fieldDicMaxType := fieldDic["StandardName"]
				fieldMetaMaxType := fieldMeta["maximum_length"]
				detail := models.LogDetails{
					ID:           bson.NewObjectId().Hex(),
					LogID:        log.ID,
					Synonymous:   cnt == 0,
					MaxLength:    fieldDicMaxLength == fieldMetaMaxLength,
					Type:         fieldDicMaxType == fieldMetaMaxType,
					FieldID:      fieldMeta["id"],
					StandardName: field,
					ResourceName: resourceName,
					SystemName:   fieldMeta["system_name"],
					ClassID:      fieldMeta["class_id"],
					ResourceID:   fieldMeta["resource_id"],
					Compare:      compare,
				}

				if err := removeMapped(mlsSession, field, fieldMeta); err != nil {
					return fmt.Errorf("1040 %v", err)
				}
				if detail.Synonymous /* && detail.MaxLength && detail.Type */ {
					detail.Status = "ALL"
				} else {
					detail.Status = "PARTIAL"
				}
				details = append(details, detail)
				if err := assignMapped(mlsSession, field, detail.FieldID, fieldMeta, detail.Status); err != nil {
					return fmt.Errorf("1045 %v", err)
				}

				dicField := orm.Model{
					"FieldID":     "",
					"FieldStatus": "",
				}
				dicField.SetTableName(resourceName)
				cond := fmt.Sprintf("FieldID = '%v'", detail.FieldID)
				if err := hostSession.Update(dicField, cond); err != nil {
					return fmt.Errorf("1055 %v", err)
				}

				dicField["FieldID"] = detail.FieldID
				dicField["FieldStatus"] = "ALL"

				exists := map[string]bool{}
				_synonym := make([]string, 0)
				for _, fieldValidate := range fieldsValidate {
					value := fieldMeta[fieldValidate]
					if _, ok := exists[value]; !ok {
						exists[value] = true
						_synonym = append(_synonym, value)
					}
				}
				dicField["Synonym"] = strings.Join(_synonym, ",")
				cond = fmt.Sprintf("StandardName = '%v'", field)
				if err := hostSession.Update(dicField, cond); err != nil {
					return fmt.Errorf("1050 %v", err)
				}
				isynonym := models.Synonym{
					ID:           bson.NewObjectId().Hex(),
					ClassID:      detail.ClassID,
					ResourceID:   detail.ResourceID,
					ResourceName: detail.ResourceName,
					StandardName: field,
					FieldID:      detail.FieldID,
					SystemName:   detail.SystemName,
					Synonymous:   dicField["Synonym"].(string),
					ClassName:    fieldMeta["class_name"],
				}
				sql := fmt.Sprintf(`DELETE FROM %v 
				                     WHERE standard_name = ? 
									   AND field_id = ?`, isynonym.TableName())
				if _, err := hostConn.Exec(sql, field, detail.FieldID); err != nil {
					return fmt.Errorf("1060 %v", err)
				}
				if _, err := hostConn.Insert(isynonym); err != nil {
					return fmt.Errorf("1070 %v", err)
				}
				break
			}
		}
		if _found == false {
			detail := models.LogDetails{}
			detail.ID = bson.NewObjectId().Hex()
			detail.LogID = log.ID
			detail.Synonymous = false
			detail.MaxLength = false
			detail.Type = false
			detail.FieldID = fieldMeta["id"]
			detail.StandardName = ""
			detail.ResourceName = resourceName
			detail.SystemName = fieldMeta["system_name"]
			detail.ClassID = fieldMeta["class_id"]
			detail.ResourceID = fieldMeta["resource_id"]
			detail.Compare = map[string]bool{}
			detail.Status = ""
			details = append(details, detail)
			removeMapped(mlsSession, detail.FieldID, fieldMeta)

			dicField := orm.Model{
				"FieldID":     "",
				"FieldStatus": "",
			}
			dicField.SetTableName(resourceName)
			cond := fmt.Sprintf("FieldID = '%v'", detail.FieldID)
			if err := hostSession.Update(dicField, cond); err != nil {
				return fmt.Errorf("1055 %v", err)
			}
		}
		if len(details) == maxInsert {
			if _, err := hostConn.Insert(details); err != nil {
				return fmt.Errorf("1080 %v", err)
			}
			details = make([]models.LogDetails, 0, maxInsert)
		}
	}
	if err != nil {
		return fmt.Errorf("1090 %v", err)
	}
	return nil
}

//GetResources Get the metadata and store it in bd
func GetResources(mlsCon *xorm.Engine, mlsID string) error {
	rows := make([]models.Mls, 0)
	helper.Sync(mlsCon)
	mls.Search(&rows, mlsID)
	if len(rows) == 0 {
		return fmt.Errorf("MLS not found")
	}
	config := common.Config{
		Username:    rows[0].Username,
		Password:    rows[0].Password,
		URL:         rows[0].URL,
		Version:     rows[0].VersionRets,
		UserAgent:   rows[0].UseragentName,
		UserAgentPw: rows[0].UseragentPassword,
	}
	if rows[0].UseragentName == "" {
		config.UserAgent = "Threewide/1.0"
	}
	meta := getMetadata(config)
	resources, err := importResources(mlsCon, mlsID, meta)
	if len(err) != 0 {
		return fmt.Errorf("1%v", err[0])
	}
	classes, err := importClass(mlsCon, mlsID, meta, resources)
	if len(err) != 0 {
		return fmt.Errorf("2%v", err[0])
	}
	_, err = importTable(mlsCon, mlsID, meta, classes)
	if len(err) != 0 {
		return fmt.Errorf("3%v", err[0])
	}
	return nil
}

func importResources(conn *xorm.Engine, id string, meta parse.Meta) (map[string]string, []error) {
	conn.Delete(models.Resource{Mls: id})
	result := map[string]string{}
	errors := make([]error, 0)
	data := make([]models.Resource, 0)
	for _, v := range meta.Rets.MetadataResource {
		for k := range v.Data {
			resource := models.Resource{
				ID:                          bson.NewObjectId().Hex(),
				Mls:                         id,
				ResourceID:                  v.Get(k, "ResourceID"),
				StandardName:                v.Get(k, "StandardName"),
				VisibleName:                 v.Get(k, "VisibleName"),
				Description:                 v.Get(k, "Description"),
				KeyField:                    v.Get(k, "KeyField"),
				ClassCount:                  v.Get(k, "ClassCount"),
				ClassVersion:                v.Get(k, "ClassVersion"),
				ClassDate:                   v.Get(k, "ClassDate"),
				ObjectVersion:               v.Get(k, "ObjectVersion"),
				ObjectDate:                  v.Get(k, "ObjectDate"),
				SearchHelpVersion:           v.Get(k, "SearchHelpVersion"),
				SearchHelpDate:              v.Get(k, "SearchHelpDate"),
				EditMaskVersion:             v.Get(k, "EditMaskVersion"),
				EditMaskDate:                v.Get(k, "EditMaskDate"),
				LookupVersion:               v.Get(k, "LookupVersion"),
				LookupDate:                  v.Get(k, "LookupDate"),
				UpdateHelpVersion:           v.Get(k, "UpdateHelpVersion"),
				UpdateHelpDate:              v.Get(k, "UpdateHelpDate"),
				ValidationExpressionVersion: v.Get(k, "ValidationExpressionVersion"),
				ValidationExpressionDate:    v.Get(k, "ValidationExpressionDate"),
				ValidationLookupVersion:     v.Get(k, "ValidationLookupVersion"),
				ValidationLookupDate:        v.Get(k, "ValidationLookupDate"),
				ValidationExternalVersion:   v.Get(k, "ValidationExternalVersion"),
				ValidationExternalDate:      v.Get(k, "ValidationExternalDate"),
			}
			data = append(data, resource)
			_, err := conn.Insert(resource)
			result[resource.ResourceID] = resource.ID
			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	return result, errors
}

func getMetadata(config common.Config) parse.Meta {
	document, _ := rets.GetMetadata(config)
	meta := parse.Meta{}
	meta.LoadData(document)
	return meta
}

//NewMapResource Creates an object capable of mapping a resource
func NewMapResource(hostConn *xorm.Engine, mlsConn *xorm.Engine, classID string) *MapResource {
	return &MapResource{
		hostConn: hostConn,
		mlsConn:  mlsConn,
		classID:  classID,
		hostOrm:  orm.NewSession(hostConn),
		mlsOrm:   orm.NewSession(mlsConn),
	}
}

//MapResource object capable of mapping a resource
type MapResource struct {
	hostConn *xorm.Engine
	mlsConn  *xorm.Engine
	hostOrm  *orm.Session
	mlsOrm   *orm.Session
	classID  string
	class    models.Class
	tables   []models.Table
}

//Add Adds the requested class to the data dictionary
func (el MapResource) Add() error {
	el.class = models.Class{}
	el.tables = make([]models.Table, 0)
	metadata.GetTables(el.mlsConn, &el.tables, el.classID)
	metadata.GetClassByID(el.mlsConn, &el.class, el.classID)
	el.createTable()
	el.mapFields()
	ClassStoreAs(el.hostConn, el.classID, el.class.ResourceName)
	if err := el.createResource(); err != nil {
		return err
	}
	resource := el.class.ResourceName
	row := models.LookupFieldsAndValues{
		LookupValue:      "ResourceName",
		LookupField:      el.class.ResourceName,
		Definition:       fmt.Sprintf("This record is related to another record in the %v resource.", resource),
		LookupStatus:     "Active",
		StatusChangeDate: "20150721T0000",
		RevisedDate:      "20150721T0000",
		AddedInVersion:   "1.4.0",
		WikiPage:         fmt.Sprintf("%v (ResourceName)", resource),
	}
	sql := fmt.Sprintf("DELETE FROM %v WHERE LookupValue = ? AND LookupField = ?", row.TableName())
	if _, err := el.hostConn.Exec(sql, row.LookupValue, row.LookupField); err != nil {
		return err
	}
	if _, err := el.hostConn.InsertOne(row); err != nil {
		return err
	}
	return nil
}

//Map Performs resource mapping
func (el MapResource) Map() error {
	el.class = models.Class{}
	el.tables = make([]models.Table, 0)
	metadata.GetTables(el.mlsConn, &el.tables, el.classID)
	metadata.GetClassByID(el.mlsConn, &el.class, el.classID)
	el.createTable()
	el.mapFields()
	ClassStoreAs(el.hostConn, el.classID, el.class.ResourceName)
	if err := el.createResource(); err != nil {
		return err
	}
	resource := el.class.ResourceName
	row := models.LookupFieldsAndValues{
		LookupValue:      "ResourceName",
		LookupField:      el.class.ResourceName,
		Definition:       fmt.Sprintf("This record is related to another record in the %v resource.", resource),
		LookupStatus:     "Active",
		StatusChangeDate: "20150721T0000",
		RevisedDate:      "20150721T0000",
		AddedInVersion:   "1.4.0",
		WikiPage:         fmt.Sprintf("%v (ResourceName)", resource),
	}
	sql := fmt.Sprintf("DELETE FROM %v WHERE LookupValue = ? AND LookupField = ?", row.TableName())
	if _, err := el.hostConn.Exec(sql, row.LookupValue, row.LookupField); err != nil {
		return err
	}
	if _, err := el.hostConn.InsertOne(row); err != nil {
		return err
	}
	return nil
}

func (el MapResource) mapFields() {
	model := models.DictionaryResource{}
	model.Table = el.class.ResourceName
	for _, v := range el.tables {
		StandardName := v.StandardName
		if StandardName == "" {
			StandardName = v.SystemName
		}
		cnt, _ := el.hostConn.
			Where("StandardName = ?", StandardName).
			Count(model)
		if cnt > 0 {
			continue
		}
		row := models.DictionaryResource{
			StandardName:       StandardName,
			Definition:         "",
			SimpleDataType:     v.DataType,
			SugMaxLength:       v.MaximumLength,
			Synonym:            "",
			ElementStatus:      "Active",
			BEDES:              "",
			CertificationLevel: "Silver",
			RecordID:           "",
			LookupStatus:       "<n/a>",
			Lookup:             "<n/a>",
			SugMaxPrecision:    "",
			RepeatingElement:   "No",
			PropertyTypes:      "",
			Payloads:           "",
			StatusChangeDate:   "",
			RevisedDate:        "",
			AddedInVersion:     "",
			Wiki:               "",
			Table:              el.class.ResourceName,
		}
		el.hostConn.Insert(row)
	}
}

func (el MapResource) createResource() error {
	model := models.LookupFieldsAndValues{}
	cond := "LookupValue = 'ResourceName' and LookupField = '" + el.class.ResourceName + "'"
	cnt, _ := el.hostConn.Where(cond).Count(model)
	if cnt != 0 {
		return nil
	}
	def := "This record is related to another record in the " + el.class.ResourceName + " resource."
	row := models.LookupFieldsAndValues{
		LookupValue:      "ResourceName",
		LookupField:      el.class.ResourceName,
		Definition:       def,
		Synonym:          "",
		BEDES:            "",
		References:       "Media,HistoryTransactional,SavedSearch",
		LookupStatus:     "Active",
		LookupID:         "",
		LookupFieldID:    "",
		StatusChangeDate: "",
		RevisedDate:      "",
		AddedInVersion:   el.class.TableVersion,
		WikiPage:         "",
	}
	_, err := el.hostConn.Insert(row)
	return err
}

func (el MapResource) createTable() error {
	dictionary := models.DictionaryResource{Table: el.class.ResourceName}
	if err := el.hostConn.Sync2(dictionary); err != nil {
		return err
	}
	return nil
}
