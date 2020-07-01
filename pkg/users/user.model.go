// package users defines the structure of user data as will exist in the db
package users

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Credentials defines the params requisite for user creation
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Model defines a user fully
type Model struct {
	UserID string `json:"userID"`
	Credentials
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}

// ValidationErrs collects errs from validating user credentials
type ValidationErrs struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Error is defined here inorder to qualify ValidationErrs as being a valid `error` type
func (v ValidationErrs) Error() string {
	var errs ValidationErrs
	validationErrs, err := formatErrs(errs)
	if err != nil {
		return err.Error()
	}

	return string(validationErrs)
}

// formatErrs structures validation errs as a jsonified string
func formatErrs(errs ValidationErrs) ([]byte, error) {
	validationErrs, marshallingErr := json.MarshalIndent(&errs, "", "   ")
	if marshallingErr != nil {
		return nil, marshallingErr
	}

	return validationErrs, nil
}

// Validate checks that user credentials are valid
func (c *Credentials) Validate() *ValidationErrs {

	/*
	* {0,7} - is too short
	* [^0-9]* - does not contain at least one digit
	* [^A-Z]* - does not contain at least one uppercase letter
	* [^a-z]* - does not contain at least one lowercase letter
	* [a-zA-Z0-9]* - contains ONLY non-special chars
	 */

	// invalidPasswordRegex := "^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$"
	return nil
}

// CreateUser registers a new user in the db
func (m *Model) CreateUser(conn *sql.DB) error {
	return ValidationErrs{}
}

// GetByID retrieves a user from the db by id
func (m *Model) GetByID(conn *sql.DB, id string) (*Model, error) {
	return nil, nil
}

// Update updates a property of a user in the db
func (u *Model) Update(conn *sql.DB, id string) error {
	return nil
}
