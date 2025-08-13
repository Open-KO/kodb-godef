package jsonSchema

import (
	"fmt"
	"github.com/Open-KO/kodb-godef/enums/dbType"
	"github.com/Open-KO/kodb-godef/enums/profile"
	"github.com/Open-KO/kodb-godef/enums/tsql"
)

// READ BEFORE MODIFYING:
// This project is used as a submodule for ko-codegen and kodb-util
// if you're modifying this, you'll need to add the following to the end of your
// importing go.mod to use your local changes.  Otherwise, you'll be using
// the main branch code from github
// replace github.com/Open-KO/OpenKO-db => ./OpenKO-db

type TableDef struct {
	Database    dbType.DbType `json:"database"`          // Which database this table should be created inside of
	Name        string        `json:"name"`              // Table name in the database
	ClassName   string        `json:"className"`         // Code-friendly class/struct name
	Description string        `json:"description"`       // Table description
	Exports     []Export      `json:"exports,omitempty"` // application specific exports
	Unions      []Union       `json:"unions,omitempty"`  // grouped columns to generate unions from
	Indexes     []IndexDef    `json:"indexes,omitempty"` // Any index definitions for the table
	Columns     []Column      `json:"columns"`           // Columns belonging to the table
}

type Column struct {
	Name          string        `json:"name"`                   // Column name in the database
	PropertyName  string        `json:"propertyName"`           // Code-friendly property name
	Description   string        `json:"description"`            // Property description
	Type          tsql.TSqlType `json:"type"`                   // Supported TSQL Type
	DefaultValue  string        `json:"defaultValue,omitempty"` // Default value that should be assigned to the property
	AllowNull     bool          `json:"allowNull,omitempty"`    // Can the column's value be null?
	Length        int           `json:"length,omitempty"`       // length specifier for array types
	Enums         []Enum        `json:"enums,omitempty"`        // array of enumerated values
	ForceBinary   bool          `json:"forceBinary,omitempty"`  // Should this column be read with a convert function to varbinary, and written back with a convert to original?
	CollationName *string       `json:"collationName,omitempty"`
	CharacterSet  *string       `json:"characterSet,omitempty"`
	AutoIncrement bool          `json:"autoIncrement,omitempty"`
}

type IndexDef struct {
	Name         string   `json:"name" gorm:"column:name"`
	Type         string   `json:"type" gorm:"column:type_desc"`
	IsUnique     bool     `json:"isUnique" gorm:"column:is_unique"`
	IsPrimaryKey bool     `json:"isPrimaryKey" gorm:"column:is_primary_key"`
	Columns      []string `json:"columns" gorm:"-"`
}

type Enum struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

type Export struct {
	Namespace profile.ExportName `json:"namespace"`         // should be the server app name
	Columns   []string           `json:"columns,omitempty"` // columns to include in the export.  If none are specified, full column export assumed
	Exclude   []string           `json:"exclude,omitempty"` // columns to exclude from the full set
}

type Union struct {
	PropertyName  string `json:"propertyName"`
	ColumnPattern string `json:"columnPattern"`
}

type ProcDef struct {
	Database    dbType.DbType `json:"database"` // Which database this proc should be executed against
	Name        string        `json:"name"`
	ClassName   string        `json:"className"`
	Description string        `json:"description"`
	Params      []ParamDef    `json:"params"`
	HasReturn   *bool         `json:"hasReturn,omitempty"`
}

type ParamDef struct {
	Name        string        `json:"name"`      // raw sql name
	ParamName   string        `json:"paramName"` // code-friendly name
	Description string        `json:"description"`
	Type        tsql.TSqlType `json:"type"`
	Length      int           `json:"length"`
	ParamIndex  int           `json:"paramIndex" gorm:"column:paramIndex"`
	IsOutput    bool          `json:"isOutput" gorm:"column:isOutput"`
	ForceBinary bool          `json:"forceBinary,omitempty"`
}

func (this *Column) GormType() string {
	if this.Length > 0 {
		return fmt.Sprintf("%s(%d)", this.Type, this.Length)
	}
	return string(this.Type)
}

func (this *Column) IsBlobType() bool {
	if this.Type == tsql.Text || this.Type == tsql.Image || this.Length > 255 {
		return true
	}
	return false
}
