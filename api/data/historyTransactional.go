package data

import "github.com/go-xorm/xorm"

func createTableHistoryTransactional(conn *xorm.Engine, tablename string) error {
	sql := `CREATE TABLE public."HistoryTransactional"
	(
	  id serial NOT NULL,
	  "HistoryTransactionalKey" text NOT NULL,
	  "OriginatingSystemHistoryKey" character varying(255) NOT NULL,
	  "OriginatingSystemName" character varying(255) NOT NULL,
	  "ChangedByMemberID" character varying(25) NOT NULL,
	  "ChangedByMemberKey" character varying(255) NOT NULL,
	  "ChangeType" character varying(255) NOT NULL,
	  "ModificationTimestamp" timestamp with time zone,
	  "FieldKey" character varying(255) NOT NULL,
	  "FieldName" character varying(255) NOT NULL,
	  "PreviousValue" character varying(16000) NOT NULL,
	  "NewValue" character varying(16000) NOT NULL,
	  "ClassName" character varying(255) NOT NULL,
	  "ResourceName" character varying(255) NOT NULL,
	  "ResourceRecordKey" character varying(255) NOT NULL,
	  "ResourceRecordID" character varying(255) NOT NULL,
	  PRIMARY KEY (id)
	)`
	_, err := conn.Exec(sql)
	return err
}
