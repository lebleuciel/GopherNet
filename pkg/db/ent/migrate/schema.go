// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// BurrowsColumns holds the columns for the "burrows" table.
	BurrowsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "depth", Type: field.TypeFloat64, Default: 0},
		{Name: "width", Type: field.TypeFloat64, Default: 0},
		{Name: "is_occupied", Type: field.TypeBool, Default: false},
		{Name: "age", Type: field.TypeInt},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// BurrowsTable holds the schema information for the "burrows" table.
	BurrowsTable = &schema.Table{
		Name:       "burrows",
		Columns:    BurrowsColumns,
		PrimaryKey: []*schema.Column{BurrowsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "burrow_name",
				Unique:  true,
				Columns: []*schema.Column{BurrowsColumns[1]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		BurrowsTable,
	}
)

func init() {
}
