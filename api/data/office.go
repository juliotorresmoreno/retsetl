package data

import "github.com/go-xorm/xorm"

func createTableOffice(conn *xorm.Engine, tablename string) error {
	sql := `CREATE TABLE public."Office"
	(
	  id serial NOT NULL,
	  "OfficeKey" text,
	  "OriginatingSystemOfficeKey" text,
	  "OriginatingSystemName" text,
	  "OfficeMlsId" text,
	  "OfficeName" text,
	  "OfficePhone" text,
	  "OfficePhoneExt" text,
	  "OfficeFax" text,
	  office_email character varying(100),
	  "OfficeType" text,
	  "OfficeBranchType" text,
	  "SocialMedia" json DEFAULT '{}'::json,
	  "OfficeAOR" text,
	  "OfficeAORMlsId" text,
	  "OfficeAORkey" text,
	  "OfficeNationalAssociationId" text,
	  "OfficeCorporateLicense" text,
	  "OfficeBrokerMlsId" text,
	  "OfficeBrokerKey" text,
	  "OfficeManagerMlsId" text,
	  "OfficeManagerKey" text,
	  "OfficeAddress1" text,
	  "OfficeAddress2" text,
	  "OfficeCity" text,
	  "OfficeStateOrProvince" text,
	  "OfficePostalCode" text,
	  "OfficePostalCodePlus4" text,
	  "OfficeCountyOrParish" text,
	  "OfficeStatus" text,
	  "OfficeAssociationComments" text,
	  "OriginalEntryTimestamp" timestamp with time zone,
	  "ModificationTimestamp" timestamp with time zone,
	  "DeletedAt" timestamp with time zone,
	  "MainOfficeKey" text,
	  "MainOfficeMlsId" text,
	  "FranchiseAffiliation" text,
	  "IDXOfficeParticipationYN" boolean,
	  "SyndicateTo" json DEFAULT '{}'::json,
	  "SyndicateAgentOption" text,
	  CONSTRAINT "Offices_pkey" PRIMARY KEY (id)
	)`
	_, err := conn.Exec(sql)
	return err
}
