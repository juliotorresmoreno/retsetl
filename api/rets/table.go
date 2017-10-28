package rets

import (
	"github.com/go-xorm/xorm"

	"gopkg.in/mgo.v2/bson"

	"bitbucket.org/mlsdatatools/retsetl/models"
	"bitbucket.org/mlsdatatools/retsetl/parse"
)

func importTable(conn *xorm.Engine, mls string, meta parse.Meta, classes map[string]map[string]string) (bool, []error) {
	conn.Delete(models.Table{Mls: mls})
	limit := 20
	result := map[string]string{}
	errors := make([]error, 0)
	tables := make([]models.Table, 0, limit)
	for _, v := range meta.Rets.MetadataTable {
		if resource, ok := classes[v.Resource]; ok {
			if class, ok := resource[v.Class]; ok {
				length := len(v.Data)
				for k := range v.Data {
					row := models.Table{
						ID:             bson.NewObjectId().Hex(),
						ClassID:        class,
						ClassName:      v.Class,
						ResourceID:     resource["id"],
						ResourceName:   v.Resource,
						Mls:            mls,
						SystemName:     v.Get(k, "SystemName"),
						StandardName:   v.Get(k, "StandardName"),
						LongName:       v.Get(k, "LongName"),
						DBName:         v.Get(k, "DBName"),
						ShortName:      v.Get(k, "ShortName"),
						MaximumLength:  v.Get(k, "MaximumLength"),
						DataType:       v.Get(k, "DataType"),
						Precision:      v.Get(k, "Precision"),
						Searchable:     v.Get(k, "Searchable"),
						Interpretation: v.Get(k, "Interpretation"),
						Alignment:      v.Get(k, "Alignment"),
						UseSeparator:   v.Get(k, "UseSeparator"),
						EditMaskID:     v.Get(k, "EditMaskID"),
						LookupName:     v.Get(k, "LookupName"),
						MaxSelect:      v.Get(k, "MaxSelect"),
						Units:          v.Get(k, "Units"),
						Index:          v.Get(k, "Index"),
						Minimum:        v.Get(k, "Minimum"),
						Maximum:        v.Get(k, "Maximum"),
						Default:        v.Get(k, "Default"),
						Required:       v.Get(k, "Required"),
						SearchHelpID:   v.Get(k, "SearchHelpID"),
						Unique:         v.Get(k, "Unique"),
					}
					tables = append(tables, row)
					result[row.ClassID] = row.ID
					if len(tables) == limit || k == length-1 {
						_, err := conn.Insert(tables)
						tables = make([]models.Table, 0, limit)
						if err != nil {
							errors = append(errors, err)
						}
					}
				}
			}
		}
	}
	return len(errors) == 0, errors
}

/*
func GetAliasField(conn *xorm.Engine, columnID string, rows *[]models.DictionaryAlias) (bool, error) {
	if err := conn.Where("column_id = ?", columnID).Find(rows); err != nil {
		return false, err
	}
	return true, nil
}*/
/*
func SetAliasField(conn *xorm.Engine, columnID, name string) (bool, error) {
	var err error
	name = strings.Replace(name, " ", ",", -1)
	alias := strings.Split(name, ",")
	rows := make([]models.DictionaryAlias, 0, len(alias))
	column := models.DictionaryColumn{}
	conn.Where("id = ?", columnID).Get(&column)
	if column.ID != columnID {
		return false, fmt.Errorf("Not found")
	}
	for _, v := range alias {
		row := models.DictionaryAlias{
			ID:       bson.NewObjectId().Hex(),
			Name:     v,
			ColumnID: columnID,
			TableID:  column.TableID,
		}
		rows = append(rows, row)
	}
	_, err = conn.Where(fmt.Sprintf("column_id = '%v'", columnID)).Delete(models.DictionaryAlias{})
	if err != nil {
		return false, err
	}
	_, err = conn.Insert(rows)
	if err != nil {
		return false, err
	}
	return true, nil
}*/
/*
func GetIndexTable(conn *xorm.Engine, tableID string, row *models.DictionaryTable) (bool, error) {
	if _, err := conn.Id(tableID).Get(row); err != nil {
		return false, err
	}
	return true, nil
}
*/
/*
func SetIndexTable(conn *xorm.Engine, tableID, name string) (bool, error) {
	var err error
	row := models.DictionaryTable{}
	if _, err = conn.Id(tableID).Get(&row); err != nil {
		return false, err
	}
	table := models.Class{}
	if _, err = conn.Id(tableID).Get(&table); err != nil {
		return false, err
	}
	if table.ID != tableID {
		return false, fmt.Errorf("Not found")
	}
	row.Name = name
	if row.ID == tableID {
		_, err = conn.Id(tableID).Update(row)
	} else {
		row.ID = tableID
		_, err = conn.Insert(row)
		indexFields(conn, tableID)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
*/
/*
func indexFields(conn *xorm.Engine, tableID string) {
	rows := make([]models.Table, 0)
	metadata.GetTables(&rows, tableID)
	fields := make([]models.DictionaryColumn, 0, len(rows))
	alias := make([]models.DictionaryAlias, 0, len(rows))
	for _, field := range rows {
		_field := models.DictionaryColumn{
			ID:           field.ID,
			Name:         strings.Trim(field.SystemName, " "),
			TableID:      tableID,
			SystemName:   field.SystemName,
			DBName:       field.DBName,
			DataType:     field.DataType,
			Searchable:   field.Searchable,
			UseSeparator: field.UseSeparator,
			Index:        field.Index,
			Default:      field.Default,
		}
		fields = append(fields, _field)
		_alias := make([]string, 0)
		_alias = append(_alias, field.SystemName)

		if ok, _ := helper.InArray(field.StandardName, _alias); field.StandardName != "" && !ok {
			alias = append(alias, models.DictionaryAlias{
				ID:       bson.NewObjectId().Hex(),
				Name:     field.StandardName,
				ColumnID: _field.ID,
				TableID:  tableID,
			})
			_alias = append(_alias, strings.Trim(field.StandardName, " "))
		}
		if ok, _ := helper.InArray(field.StandardName, _alias); field.LongName != "" && !ok {
			alias = append(alias, models.DictionaryAlias{
				ID:       bson.NewObjectId().Hex(),
				Name:     field.LongName,
				ColumnID: _field.ID,
				TableID:  tableID,
			})
			_alias = append(_alias, strings.Trim(field.LongName, " "))
		}
		if ok, _ := helper.InArray(field.StandardName, _alias); field.ShortName != "" && !ok {
			alias = append(alias, models.DictionaryAlias{
				ID:       bson.NewObjectId().Hex(),
				Name:     field.ShortName,
				ColumnID: _field.ID,
				TableID:  tableID,
			})
		}
	}
	conn.Insert(fields)
	conn.Insert(alias)
}
*/
