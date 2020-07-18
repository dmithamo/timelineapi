package models

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/dmithamo/timelineapi/pkg/dbservice"
	"github.com/dmithamo/timelineapi/pkg/security"
	"github.com/dmithamo/timelineapi/pkg/utils"
)

// UserCredentials defines the params requisite for user creation
type UserCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// User defines a user fully
type User struct {
	UserID string `json:"userID,omitempty"`
	UserCredentials
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
func (c UserCredentials) Error() string {
	return "err in user credentials"
}

// Validate checks that user credentials are valid
func (c *UserCredentials) Validate() error {

	validationErrs := &UserCredentials{}
	hasErrors := false

	func() {
		if invalidPasswordRegex.MatchString(c.Password) {
			if c.Password == "" {
				validationErrs.Password = "password is required"
			} else {
				validationErrs.Password = invalidPasswordMessage
			}
			hasErrors = true
		}
	}()

	func() {
		if !validEmailRegex.MatchString(c.Username) {
			if c.Username == "" {
				validationErrs.Username = "username is required"
			} else {
				validationErrs.Username = invalidEmailMessage
			}
			hasErrors = true
		}
	}()

	if !hasErrors {
		return nil
	}

	return validationErrs
}

// CreateUser registers a new user in the db
func (u *User) CreateUser(db *sql.DB, credentials *UserCredentials) error {
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
func (u *User) GetByCredentials(db *sql.DB, credentials *UserCredentials) (*User, error) {
	var pwdHash string
	var userID string
	var user User

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

// GetByUUID searches the db for a user with a given UUID
func (u *User) GetByUUID(db *sql.DB, uuid string) (*User, error) {
	stmt, err := db.Prepare("SELECT BIN_TO_UUID(userID) userID FROM users WHERE userID=UUID_TO_BIN(?)")
	if err != nil {
		return nil, err
	}

	var user *User
	err = stmt.QueryRow(uuid).Scan(user.UserID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdatePassword updates an action's title or description
func (u *User) UpdatePassword(db *sql.DB, userID string, password string) error {
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
