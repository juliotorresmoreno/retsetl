package models

//Table Represents the table in the metadata of a server rets
type Table struct {
	ID             string `xorm:"id varchar(40) not null pk"            json:"id"`
	ClassID        string `xorm:"'class_id' varchar(255) not null"      json:"class_id"`
	ClassName      string `xorm:"'class_name' varchar(255) not null"    json:"class_name"`
	ResourceID     string `xorm:"'resource_id' varchar(255) not null"   json:"resource_id"`
	ResourceName   string `xorm:"'resource_name' varchar(255) not null" json:"resource_name"`
	Mls            string `xorm:"mls varchar(255) not null"             json:"mls"`
	Mapped         string `xorm:"mapped varchar(255) not null"          json:"mapped"`
	MappedStatus   string `xorm:"mapped_status varchar(255) not null"   json:"mapped_status"`
	SystemName     string `xorm:"system_name varchar(255) not null"     json:"system_name"`
	StandardName   string `xorm:"standard_name varchar(255) not null"   json:"standard_name"`
	LongName       string `xorm:"long_name varchar(255) not null"       json:"long_name"`
	DBName         string `xorm:"db_name varchar(255) not null"         json:"db_name"`
	ShortName      string `xorm:"short_name varchar(255) not null"      json:"short_name"`
	MaximumLength  string `xorm:"maximum_length varchar(255) not null"  json:"maximum_length"`
	DataType       string `xorm:"data_type varchar(255) not null"       json:"data_type"`
	Precision      string `xorm:"precision varchar(255) not null"       json:"precision"`
	Searchable     string `xorm:"searchable varchar(255) not null"      json:"searchable"`
	Interpretation string `xorm:"interpretation varchar(255) not null"  json:"interpretation"`
	Alignment      string `xorm:"alignment varchar(255) not null"       json:"alignment"`
	UseSeparator   string `xorm:"use_separator varchar(255) not null"   json:"use_separator"`
	EditMaskID     string `xorm:"edit_mask_id varchar(255) not null"    json:"edit_mask_id"`
	LookupName     string `xorm:"lookup_name varchar(255) not null"     json:"lookup_name"`
	MaxSelect      string `xorm:"max_select varchar(255) not null"      json:"max_select"`
	Units          string `xorm:"units varchar(255) not null"           json:"units"`
	Index          string `xorm:"'index' varchar(255) not null"         json:"index"`
	Minimum        string `xorm:"minimum varchar(255) not null"         json:"minimum"`
	Maximum        string `xorm:"maximum varchar(255) not null"         json:"maximum"`
	Default        string `xorm:"'default' varchar(255) not null"       json:"default"`
	Required       string `xorm:"required varchar(255) not null"        json:"required"`
	SearchHelpID   string `xorm:"search_help_id varchar(255) not null"  json:"search_help_id"`
	Unique         string `xorm:"'unique' varchar(255) not null"        json:"unique"`
}

//TableName Name of the table where the data will be stored
func (el Table) TableName() string {
	return "metadata_tables"
}
