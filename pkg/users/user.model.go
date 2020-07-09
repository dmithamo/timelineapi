// package users defines the structure of user data as will exist in the db
package users

import (
	"database/sql"
	"encoding/json"
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

// ValidationErrs collects errs from validating user credentials
type ValidationErrs struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Error is defined here inorder to qualify ValidationErrs as being a valid `error` type
func (v ValidationErrs) Error() string {
	var errs ValidationErrs
	validationErrs, err := FormatErrs(errs)
	if err != nil {
		return err.Error()
	}

	return string(validationErrs)
}

// FormatErrs structures validation errs as a jsonified string
func FormatErrs(errs ValidationErrs) ([]byte, error) {
	validationErrs, marshallingErr := json.MarshalIndent(&errs, "", "   ")
	if marshallingErr != nil {
		return nil, marshallingErr
	}

	return validationErrs, nil
}

// Validate checks that user credentials are valid
func (c *Credentials) Validate() error {
	validEmailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	invalidPasswordRegex := regexp.MustCompile("^(.{0,7}|[^0-9]*|[^A-Z]*|[^a-z]*|[a-zA-Z0-9]*)$")

	isInvalidPassword := invalidPasswordRegex.MatchString(c.Password)
	isvalidUsername := validEmailRegex.MatchString(c.Username)

	hasErrors := false
	validationErrs := &ValidationErrs{}
	if isInvalidPassword {
		validationErrs.Password = "password should be at least 8 characters long, combining uppercase letters, lowercase letters, digits, and special characters"
		hasErrors = true
	}

	if !isvalidUsername {
		validationErrs.Username = "username should be a valid email address"
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

// Update updates a property of a user in the db
func (u *Model) Update(db *sql.DB, id string) error {
	return nil
}
