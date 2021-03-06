package data

import "github.com/go-xorm/xorm"

func createTableContact(conn *xorm.Engine, tablename string) error {
	sql := `CREATE TABLE public."Contact"
	(
	  id serial NOT NULL,
	  "ContactKey" character varying(255),
	  "ContactLoginId" character varying(25),
	  "ContactPassword" character varying(25),
	  "OriginatingSystemContactKey" character varying(255),
	  "OriginatingSystemName" character varying(255),
	  "OwnerMemberID" character varying(25),
	  "NamePrefix" character varying(10),
	  "FirstName" character varying(50),
	  "MiddleName" character varying(50),
	  "LastName" character varying(50),
	  "NameSuffix" character varying(10),
	  "FullName" character varying(150),
	  "Nickname" character varying(50),
	  "ReferredBy" character varying(150),
	  "JobTitle" character varying(50),
	  "Notes" character varying(1024),
	  "HomeAddress1" character varying(50),
	  "HomeAddress2" character varying(50),
	  "HomeCity" character varying(50),
	  "HomeStateOrProvince" json DEFAULT '{}'::json,
	  "HomePostalCode" character varying(10),
	  "HomePostalCodePlus4" character varying(4),
	  "HomeCarrierRoute" character varying(9),
	  "HomeCountyOrParish" json DEFAULT '{}'::json,
	  "HomeCountry" json DEFAULT '{}'::json,
	  "WorkAddress1" character varying(50),
	  "WorkAddress2" character varying(50),
	  "WorkCity" character varying(50),
	  "WorkStateOrProvince" character varying(2),
	  "WorkPostalCode" character varying(10),
	  "WorkPostalCodePlus4" character varying(4),
	  "WorkCarrierRoute" character varying(9),
	  "WorkCountyOrParish" character varying(50),
	  "WorkCountry" character varying(2),
	  "OtherAddress1" character varying(50),
	  "OtherAddress2" character varying(50),
	  "OtherCity" character varying(50),
	  "OtherStateOrProvince" character varying(2),
	  "OtherPostalCode" character varying(10),
	  "OtherPostalCodePlus4" character varying(4),
	  "OtherCarrierRoute" character varying(9),
	  "OtherCountyOrParish" character varying(50),
	  "OtherCountry" character varying(2),
	  "PreferredAddress" character varying(255),
	  "PreferredPhone" character varying(255),
	  "Email" character varying(80),
	  "Email2" character varying(80),
	  "Email3" character varying(80),
	  "OfficePhone" character varying(16),
	  "OfficePhoneExt" integer,
	  "MobilePhone" character varying(16),
	  "DirectPhone" character varying(16),
	  "HomePhone" character varying(16),
	  "HomeFax" character varying(16),
	  "BusinessFax" character varying(16),
	  "Pager" character varying(16),
	  "VoiceMail" character varying(16),
	  "VoiceMailExt" integer,
	  "TollFreePhone" character varying(16),
	  "PhoneTTYTTD" character varying(16),
	  "OtherPhoneType" character varying(25),
	  "OtherPhoneTypeNumber" json DEFAULT '{}'::json,
	  "OtherPhoneTypeExt" json DEFAULT '{}'::json,
	  "Company" character varying(50),
	  "Department" character varying(50),
	  "SocialMediaType" character varying(25),
	  "SocialMediaTypeUrlOrId" json DEFAULT '{}'::json,
	  "Birthdate" timestamp with time zone,
	  "Anniversary" timestamp with time zone,
	  "OriginalEntryTimestamp" timestamp with time zone,
	  "ModificationTimestamp" timestamp with time zone,
	  "DeletedAt" timestamp with time zone,
	  "UserDefinedFieldName" json DEFAULT '{}'::json,
	  "UserDefinedFieldValue" json DEFAULT '{}'::json,
	  "AssistantName" character varying(150),
	  "AssistantPhone" character varying(16),
	  "AssistantPhoneExt" integer,
	  "AssistantEmail" character varying(80),
	  "SpousePartnerName" character varying(150),
	  "Children" character varying(150),
	  "Gender" character varying(6),
	  "Language" json DEFAULT '{}'::json,
	  "Groups" json DEFAULT '{}'::json,
	  "ContactStatus" character varying(25),
	  "ContactType" json DEFAULT '{}'::json,
	  contact_key text,
	  PRIMARY KEY (id)
	)`
	_, err := conn.Exec(sql)
	return err
}
