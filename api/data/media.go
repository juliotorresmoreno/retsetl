package data

import "github.com/go-xorm/xorm"

func createTableMedia(conn *xorm.Engine, tablename string) error {
	sql := `CREATE TABLE public."Media"
	(
	  id serial NOT NULL,
	  "MediaKey" character varying(255),
	  "ResourceRecordKey" character varying(255),
	  "ResourceRecordID" character varying(255),
	  "OriginatingSystemMediaKey" character varying(255),
	  "OriginatingSystemName" character varying(255),
	  "MediaObjectID" character varying(255),
	  "ChangedByMemberID" character varying(255),
	  "ChangedByMemberKey" character varying(255),
	  "MediaCategory" character varying(50),
	  "MimeType" character varying(50),
	  "ShortDescription" character varying(50),
	  "LongDescription" character varying(1024),
	  "OriginalEntryTimestamp" timestamp with time zone,
	  "ModificationTimeStamp" timestamp with time zone,
	  "DeletedAtTimestamp" timestamp with time zone,
	  "MediaModificationTimestamp" timestamp with time zone,
	  "MediaURL" character varying(8000),
	  "MediaHTML" character varying(8500),
	  "Order" integer,
	  "Group" character varying(50),
	  "ImageWidth" integer,
	  "ImageHeight" integer,
	  "ImageSizeDescription" character varying(50),
	  "ResourceName" character varying(50),
	  "ClassName" character varying(50),
	  "Permission" json DEFAULT '{}'::json,
	  "MediaStatus" character varying(50),
	  PRIMARY KEY (id)
	)`
	_, err := conn.Exec(sql)
	return err
}
