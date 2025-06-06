// Code generated by ent, DO NOT EDIT.

package burrow

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the burrow type in the database.
	Label = "burrow"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDepth holds the string denoting the depth field in the database.
	FieldDepth = "depth"
	// FieldWidth holds the string denoting the width field in the database.
	FieldWidth = "width"
	// FieldIsOccupied holds the string denoting the is_occupied field in the database.
	FieldIsOccupied = "is_occupied"
	// FieldAge holds the string denoting the age field in the database.
	FieldAge = "age"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the burrow in the database.
	Table = "burrows"
)

// Columns holds all SQL columns for burrow fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDepth,
	FieldWidth,
	FieldIsOccupied,
	FieldAge,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultDepth holds the default value on creation for the "depth" field.
	DefaultDepth float64
	// DefaultWidth holds the default value on creation for the "width" field.
	DefaultWidth float64
	// DefaultIsOccupied holds the default value on creation for the "is_occupied" field.
	DefaultIsOccupied bool
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(int) error
)

// OrderOption defines the ordering options for the Burrow queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDepth orders the results by the depth field.
func ByDepth(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDepth, opts...).ToFunc()
}

// ByWidth orders the results by the width field.
func ByWidth(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWidth, opts...).ToFunc()
}

// ByIsOccupied orders the results by the is_occupied field.
func ByIsOccupied(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsOccupied, opts...).ToFunc()
}

// ByAge orders the results by the age field.
func ByAge(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAge, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
