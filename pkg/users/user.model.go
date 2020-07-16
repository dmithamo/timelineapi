// package users contains functionality for managing user data
package users

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/dmithamo/timelineapi/pkg/dbservice"
	"github.com/dmithamo/timelineapi/pkg/security"
	"github.com/dmithamo/timelineapi/pkg/utils"
)

// Credentials defines the params requisite for user creation
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Model defines a user fully
type Model struct {
	UserID string `json:"userID,omitempty"`
	Credentials
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// regexes for valid creds
var validEmailRegex *regexp.Regexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var invalidPasswordRegex *regexp.Regexp = regexp.MustCompile("^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$")

// error messages
var invalidEmailMessage string = "invalid username. Use a valid email address"
var invalidPasswordMessage string = "invalid password. Use at least 8 characters, combining uppercase letters, lowercase letters, digits, and special characters"

// Error is defined here inorder to qualify validation errs (modelled on Credentials)
// as being a valid `error` type
// It is otherwise quite useless
func (c Credentials) Error() string {
	return "err in user credentials"
}

// Validate checks that user credentials are valid
func (c *Credentials) Validate() error {

	hasErrors := false
	validationErrs := &Credentials{}

	if invalidPasswordRegex.MatchString(c.Password) {
		if c.Password == "" {
			validationErrs.Password = "password is required"
		} else {
			validationErrs.Password = invalidPasswordMessage
		}
		hasErrors = true
	}

	if !validEmailRegex.MatchString(c.Username) {
		if c.Username == "" {
			validationErrs.Username = "username is required"
		} else {
			validationErrs.Username = invalidEmailMessage
		}

		hasErrors = true
	}

	if !hasErrors {
		return nil
	}

	return validationErrs
}

// CreateUser registers a new user in the db
func (m *Model) CreateUser(db *sql.DB, credentials *Credentials) error {
	stmt, err := db.Prepare("INSERT INTO users(userID,username,password) VALUES (UUID_TO_BIN(UUID()),?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	pwdHash, err := security.GeneratePasswordHash(&credentials.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err.Error())
	}

	_, err = stmt.Exec(credentials.Username, pwdHash)

	return dbservice.CheckDatabaseErr(err, "username")
}

// GetByCredentials retrieves a user from the db by username, password - for login
func (m *Model) GetByCredentials(db *sql.DB, credentials *Credentials) (*Model, error) {
	var pwdHash string
	var userID string
	var user Model

	stmt, err := db.Prepare("SELECT BIN_TO_UUID(userID) userID, password FROM users WHERE username=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(credentials.Username).Scan(&userID, &pwdHash)
	if err != nil {
		return nil, err
	}

	isCorrectPwd := security.VerifyPassword(&pwdHash, &credentials.Password)
	if !isCorrectPwd {
		return nil, fmt.Errorf(utils.WRONG_PASSWORD_ERR)
	}

	user.UserID = userID
	return &user, nil
}

// UpdatePassword updates an action's title or description
func (m *Model) UpdatePassword(db *sql.DB, userID string, password string) error {
	if invalidPasswordRegex.MatchString(password) {
		return fmt.Errorf(invalidPasswordMessage)
	}

	stmt, err := db.Prepare("UPDATE users SET password = ? WHERE userID = ?")
	if err != nil {
		return err
	}

	pwdHash, err := security.GeneratePasswordHash(&password)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pwdHash, userID)
	if err != nil {
		return err
	}

	return nil
}
