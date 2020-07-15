// package actions contains functionality for managing actions data
package actions

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dmithamo/timelineapi/pkg/dbservice"
)

// Params defines the structure of a valid action
type Params struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Error allows for Params to be used a valid err type
func (p Params) Error() string {
	return "err in action params"
}

// Model is the interface for CRUD'ing action data in the db
type Model struct {
	ActionID string `json:"actionID,omitempty"`
	Params
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

//regexes for valid input
var validTitle = regexp.MustCompile(`^[a-zA-Z].([A-Za-z0-9_ ]){3,50}$`)
var validDescription = regexp.MustCompile(`^[a-zA-Z].([A-Za-z0-9_ ]){3,100}$`)

// Validate checks the action params for errs
func (p *Params) Validate() error {

	hasErrors := false
	validationErrs := &Params{}

	if !validTitle.MatchString(p.Title) {
		if p.Title == "" {
			validationErrs.Title = "title is required"
		} else {
			validationErrs.Title = "invalid title. Use letters, numbers and underscores only, and keep it between 4 and 50 chars long"
		}
		hasErrors = true
	}

	if !validDescription.MatchString(p.Description) {
		if p.Description == "" {
			validationErrs.Description = "description is required"
		} else {
			validationErrs.Description = "invalid description. Use letters, numbers and underscores only, and keep it between 4 and 200 chars long"
		}

		hasErrors = true
	}

	if !hasErrors {
		return nil
	}

	return validationErrs
}

// CreateAction adds a new action in the db
func (m *Model) CreateAction(db *sql.DB, params Params) error {
	stmt, err := db.Prepare("INSERT INTO actions (actionID, title, description) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(makeSlugHelper(params.Title), params.Title, params.Description)
	if err != nil {
		return dbservice.CheckDatabaseErr(err, "title")
	}

	return nil
}

// GetActions retrieves all actions in the db
func (m *Model) GetActions(db *sql.DB) ([]Model, error) {
	stmt, err := db.Prepare("SELECT * FROM actions")
	if err != nil {
		return nil, dbservice.CheckDatabaseErr(err)
	}

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	var actions []Model
	for rows.Next() {
		var action Model

		err := rows.Scan(&action.ActionID, &action.Title, &action.Description, &action.CreatedAt, &action.UpdatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// GetActionByID retrueves a single action by its actionID
func (m *Model) GetActionByID(db *sql.DB, actionID string) (*Model, error) {
	stmt, err := db.Prepare(fmt.Sprintf("SELECT * FROM actions WHERE actionID='%v'", actionID))
	if err != nil {
		return nil, err
	}

	action := Model{}
	err = stmt.QueryRow().Scan(&action.ActionID, &action.Title, &action.Description, &action.CreatedAt, &action.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

// UpdateAction updates an action's title or description
func (m *Model) UpdateAction(db *sql.DB, actionID string, params Params) error {
	updateCommand := ""
	switch {
	case params.Title != "" && params.Description != "":
		updateCommand = fmt.Sprintf("UPDATE actions SET title = %v, description = %v WHERE actionID = ?", params.Title, params.Description)

	case params.Title == "":
		updateCommand = fmt.Sprintf("UPDATE actions SET description = %v WHERE actionID = ?", params.Description)

	case params.Description == "":
		updateCommand = fmt.Sprintf("UPDATE actions SET title = %v WHERE actionID = ?", params.Title)

	default:
		return fmt.Errorf("no valid title or description in update request")

	}

	stmt, err := db.Prepare(updateCommand)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(actionID)
	if err != nil {
		return err
	}

	return nil
}

// makeSlugHelper prepares a slug given the title
func makeSlugHelper(title string) string {
	return strings.ToLower(strings.Join(strings.Split(title, " "), "-"))
}
