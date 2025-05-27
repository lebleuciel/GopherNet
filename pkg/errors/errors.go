package errors

import "fmt"

// Error types
var (
	ErrBurrowNotFound    = NewUserError("Burrow not found")
	ErrBurrowOccupied    = NewUserError("Burrow is already occupied")
	ErrBurrowNotOccupied = NewUserError("Burrow is not occupied")
	ErrInvalidBurrowID   = NewUserError("Invalid burrow ID")
	ErrDatabaseOperation = NewUserError("Database operation failed")
	ErrInternalServer    = NewUserError("Internal server error")
)

type UserError struct {
	message string
}

// NewUserError creates a new user-friendly error
func NewUserError(message string) *UserError {
	return &UserError{message: message}
}

// Error implements the error interface
func (e *UserError) Error() string {
	return e.message
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
