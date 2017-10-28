package orm

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

//Session Object that represents the connection in the database and adds additional
//functions like creation of tables and insert registers
type Session struct {
	engine *xorm.Engine
}

//Engine
func (el *Session) Engine() *xorm.Engine {
	return el.engine
}

//Close Close the session of the database
func (el *Session) Close() {
	el.engine.Close()
}

//NewSession Create a session in the database that allows to perform the operations of this orm
func NewSession(conn *xorm.Engine) *Session {
	return &Session{
		engine: conn,
	}
}

//Delete Deletes a record in the database of the specified table in the model
func (el *Session) Delete(row Model, cond string) error {
	sql := fmt.Sprintf("DELETE FROM %v WHERE %v;", row.TableName(), cond)
	BD := el.engine.DB()
	_, err := BD.Exec(sql)
	return err
}

//Update Updates a record in the database of the specified table in the model
func (el *Session) Update(row Model, cond string) error {
	values := make([]string, 0)
	sep := quotes["postgres"]
	for k, v := range row {
		if k != "tablename" {
			values = append(values, fmt.Sprintf("%v = %v%v%v", k, sep, v, sep))
		}
	}
	fieldValues := strings.Join(values, ", ")
	sql := fmt.Sprintf("UPDATE %v SET %v WHERE %v;", row.TableName(), fieldValues, cond)
	BD := el.engine.DB()
	_, err := BD.Exec(sql)
	return err
}

//Insert Insert a record in the database of the specified table in the model
func (el Session) Insert(row Model) error {
	fields := make([]string, 0)
	values := make([]string, 0)
	sep := quotes["postgres"]
	for k, v := range row {
		if k != "tablename" {
			fields = append(fields, fmt.Sprintf(`"%v"`, k))
			values = append(values, fmt.Sprintf("%v%v%v", sep, v, sep))
		}
	}
	fieldNames := strings.Join(fields, ", ")
	fieldValues := strings.Join(values, ", ")

	sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v);", row.TableName(), fieldNames, fieldValues)
	BD := el.engine.DB()
	_, err := BD.Exec(sql)
	return err
}

func mysqlRealEscapeString(value string) string {
	replace := map[string]string{
		"\\":   "\\\\",
		"'":    "&#39;",
		"\\0":  "\\\\0",
		"\n":   "\\n",
		"\r":   "\\r",
		`"`:    `\"`,
		"\x1a": "\\Z",
	}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

//InsertAll Insert a record in the database of the specified table in the model
func (el Session) InsertAll(rows []Model) error {
	if len(rows[0]) == 0 {
		return nil
	}
	conn := el.engine.DB()
	fields := make([]string, 0)
	fjson := map[string]bool{}
	if el.engine.DriverName() == "postgres" {
		sql := fmt.Sprintf(`
		    SELECT column_name, data_type
			  FROM information_schema.columns
			 WHERE table_catalog = current_database()
			   AND table_schema = 'public'
			   AND table_name = '%v'
			   AND data_type = 'json'
			`,
			rows[0].TableName(),
		)
		_result, _ := el.engine.QueryString(sql)
		for _, c := range _result {
			fjson[c["column_name"]] = true
		}
	}
	for k := range rows[0] {
		if k != "tablename" {
			fields = append(fields, fmt.Sprintf("%v", k))
		}
	}
	fieldsLength := len(fields)

	for _, row := range rows {
		_fields := make([]string, 0)
		_values := make([]string, 0)
		for i := 0; i < fieldsLength; i++ {
			value := mysqlRealEscapeString(row[fields[i]].(string))
			if value != "" {
				if _, ok := fjson[fields[i]]; !ok {
					_fields = append(_fields, "\""+fields[i]+"\"")
					_values = append(_values, "'"+value+"'")
				} else {
					_tmp, _ := json.Marshal(value)
					_fields = append(_fields, "\""+fields[i]+"\"")
					_values = append(_values, "'"+string(_tmp)+"'")
				}
			}
		}
		values := strings.Join(_values, ", ")
		fields := strings.Join(_fields, ", ")
		table := row.TableName()
		sql := fmt.Sprintf(`INSERT INTO "%v" (%v) VALUES (%v);`, table, fields, values)
		_, err := conn.Exec(sql)
		if err != nil {
			fmt.Printf("%v:\n %v\n\n", err, sql)
		}
	}
	return nil
}

//Find Allows querying the database on the table specified in the model
func (el Session) Find(row Model, cond map[string]interface{}) ([]map[string]interface{}, error) {
	values := make([]interface{}, 0)
	build := builder.Select("*").
		From(row.TableName()).
		Where(builder.Eq{"1": "1"})
	values = append(values, "1")
	other := ""
	for k, v := range cond {
		if v != nil && fmt.Sprintf("%v", v) != "" {
			switch v.(type) {
			case string:
				value := v.(string)
				if strings.Contains(value, ":") {
					switch {
					case value[0:4] == "neq:":
						build = build.Where(builder.Neq{k: "1"})
						values = append(values, value[4:])
					case value[0:5] == "like:":
						build = build.Where(builder.Like{k, "1"})
						values = append(values, " "+value[5:])
					case value[0:6] == "other:":
						other = other + " " + value[6:]
					case value[0:7] == "regexp:":
						other = other + fmt.Sprintf(" AND %v regexp (\"%v\")", k, value[7:])
					default:
						build = build.Where(builder.Eq{k: "1"})
						values = append(values, v)
					}
				} else {
					build = build.Where(builder.Eq{k: "1"})
					values = append(values, v)
				}
			case []interface{}:
				build = build.Where(builder.In(k, v.([]interface{})...))
				values = append(values, v.([]interface{})...)
			default:
				build = build.Where(builder.Eq{k: "1"})
				values = append(values, v)
			}
		}
	}
	sql, _, _ := build.ToSQL()
	sql = sql + other
	if _, ok := row["id"]; ok {
		sql = sql + " order by id asc"
	}
	var result []map[string][]byte
	var err error
	if el.engine.DriverName() == "sqlite3" {
		sql = strings.Replace(sql, "?", values[0].(string), 1)
		for i := 1; i < len(values); i++ {
			sql = strings.Replace(sql, "?", "'"+values[i].(string)+"'", 1)
		}
		result, err = el.engine.Query(sql)
	} else {
		result, err = el.engine.Query(sql, values...)
	}
	response := make([]map[string]interface{}, 0, len(result))
	for _, record := range result {
		value := map[string]interface{}{}
		for k, v := range record {
			value[k] = string(v)
		}
		response = append(response, value)
	}
	return response, err
}

//AddColumn Adds a new column to the table specified in the model
func (el Session) AddColumn(row Model, column, datatype string) error {
	columnType, err := getType("text")
	if err != nil {
		return err
	}
	sql := ""
	switch el.engine.DriverName() {
	case "mysql":
		if datatype == "" {
			sql = fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v;", row.TableName(), column, columnType)
		} else {
			sql = fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v;", row.TableName(), column, datatype)
		}
	case "postgres":
		if datatype == "" {
			sql = fmt.Sprintf(
				`ALTER TABLE "%v" ADD COLUMN "%v" %v;`,
				row.TableName(),
				column,
				columnType,
			)
		} else {
			sql = fmt.Sprintf(
				`ALTER TABLE "%v" ADD COLUMN "%v" %v;`,
				row.TableName(),
				column,
				datatype,
			)
		}
	case "sqlite3":
		if datatype == "" {
			sql = fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v;", row.TableName(), column, columnType)
		} else {
			sql = fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v;", row.TableName(), column, datatype)
		}
	case "oracle":
		sql = ""
	case "sqlsvr":
		sql = ""
	}
	_, err = el.engine.Exec(sql)
	if err != nil && !strings.Contains(err.Error(), "exists") {
		return err
	}
	return nil
}

//CreateTable Create a table on the database that will only have the id field,
//the rest of fields are added manually
func (el Session) CreateTable(row Model, driver string) error {
	if row.TableName() == "" {
		return fmt.Errorf("Table not found")
	}
	sql := ""
	switch driver {
	case "mysql":
		sql = "CREATE TABLE `" + row.TableName() + "` ("
		sql += "  \"created_at\" datetime,"
		sql += "  \"updated_at\" datetime"
		sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	case "postgres":
		sql = "CREATE TABLE \"" + row.TableName() + "\" ("
		sql += "  \"created_at\" time,"
		sql += "  \"updated_at\" time"
		sql += ");"
	case "sqlite3":
		sql = "CREATE TABLE \"" + row.TableName() + "\" ("
		sql += "  \"created_at\" datetime,"
		sql += "  \"updated_at\" datetime"
		sql += ");"
	case "oracle":
		sql = ""
	case "sqlsvr":
		sql = ""
	}
	_, err := el.engine.Exec(sql)
	return err
}
