package data

import "github.com/go-xorm/xorm"

func createTableOpenHouse(conn *xorm.Engine, tablename string) error {
	sql := `CREATE TABLE public."OpenHouse"
	(
	  id serial NOT NULL,
	  "OpenHouseKey" character varying(255) NOT NULL,
	  "OpenHouseId" character varying(255) NOT NULL,
	  "OriginatingSystemKey" character varying(255) NOT NULL,
	  "OriginatingSystemName" character varying(255) NOT NULL,
	  "ListingKey" character varying(255) NOT NULL,
	  "ListingId" character varying(255) NOT NULL,
	  "ModificationTimestamp" timestamp with time zone,
	  "OriginalEntryTimestamp" timestamp with time zone,
	  "OpenHouseDate" timestamp with time zone,
	  "OpenHouseStartTime" timestamp with time zone,
	  "OpenHouseEndTime" timestamp with time zone,
	  "ShowingAgentMlsID" character varying(25) NOT NULL,
	  "ShowingAgentKey" character varying(255) NOT NULL,
	  "ShowingAgentFirstName" character varying(50) NOT NULL,
	  "ShowingAgentLastName" character varying(50) NOT NULL,
	  "OpenHouseType" character varying(25) NOT NULL,
	  "AppointmentRequiredYN" boolean,
	  "Refreshments" character varying(255) NOT NULL,
	  "Attended" character varying(25) NOT NULL,
	  "OpenHouseRemarks" character varying(500) NOT NULL,
	  "Status" character varying(25) NOT NULL,
	  "DeletedAt" timestamp with time zone,
	  CONSTRAINT "OpenHouses_pkey" PRIMARY KEY (id)
	)`
	_, err := conn.Exec(sql)
	return err
}
