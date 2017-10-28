

PRAGMA foreign_keys = 0;

CREATE TABLE sqlitestudio_temp_table AS SELECT * FROM Contacts;

DROP TABLE Contacts;

CREATE TABLE Contacts (
    id                 INTEGER       PRIMARY KEY
                                     NOT NULL,
    StandardName       TEXT          DEFAULT NULL,
    FieldDefinition    TEXT DEFAULT NULL,
    Groups             TEXT  DEFAULT NULL,
    SimpleDataType     TEXT  DEFAULT NULL,
    SugMaxLength       INTEGER       DEFAULT NULL,
    Synonym            TEXT  DEFAULT NULL,
    ElementStatus      TEXT   DEFAULT NULL,
    BEDES              TEXT  DEFAULT NULL,
    CertificationLevel TEXT   DEFAULT NULL,
    RecordID           INTEGER       DEFAULT NULL,
    LookupStatus       TEXT  DEFAULT NULL,
    Lookup             TEXT  DEFAULT NULL,
    SugMaxPrecision    TEXT  DEFAULT NULL,
    RepeatingElement   TEXT   DEFAULT NULL,
    Payloads           TEXT  DEFAULT NULL,
    StatusChangeDate   TEXT  DEFAULT NULL,
    RevisedDate        TEXT  DEFAULT NULL,
    AddedInVersion     TEXT   DEFAULT NULL,
    WikiPage           TEXT  DEFAULT NULL,
    FieldID            TEXT,
    FieldStatus        TEXT,
    Definition         TEXT,
    PropertyTypes      TEXT,
    Wiki               TEXT
);

INSERT INTO Contacts (
                         id,
                         StandardName,
                         FieldDefinition,
                         Groups,
                         SimpleDataType,
                         SugMaxLength,
                         Synonym,
                         ElementStatus,
                         BEDES,
                         CertificationLevel,
                         RecordID,
                         LookupStatus,
                         Lookup,
                         SugMaxPrecision,
                         RepeatingElement,
                         Payloads,
                         StatusChangeDate,
                         RevisedDate,
                         AddedInVersion,
                         WikiPage,
                         FieldID,
                         FieldStatus,
                         Definition,
                         PropertyTypes,
                         Wiki
                     )
                     SELECT id,
                            StandardName,
                            FieldDefinition,
                            Groups,
                            SimpleDataType,
                            SugMaxLength,
                            Synonym,
                            ElementStatus,
                            BEDES,
                            CertificationLevel,
                            RecordID,
                            LookupStatus,
                            Lookup,
                            SugMaxPrecision,
                            RepeatingElement,
                            Payloads,
                            StatusChangeDate,
                            RevisedDate,
                            AddedInVersion,
                            WikiPage,
                            FieldID,
                            FieldStatus,
                            Definition,
                            PropertyTypes,
                            Wiki
                       FROM sqlitestudio_temp_table;

DROP TABLE sqlitestudio_temp_table;

CREATE UNIQUE INDEX UQE_Contacts_StandardName ON Contacts (
    StandardName
);

PRAGMA foreign_keys = 1;
