// package models contains functionality for managing db data
package models

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/dmithamo/timelineapi/pkg/dbservice"
)

// Params defines the structure of a valid action
type ActionParams struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Error allows for Params to be used a valid err type
func (p ActionParams) Error() string {
	return "err in action params"
}

// Action is the interface for CRUD'ing action data in the db
type Action struct {
	ActionID string `json:"actionID,omitempty"`
	ActionParams
	isArchived bool
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	UserID     string    `json:"userID,omitempty"`
}

//regexes for valid input
var validTitle = regexp.MustCompile(`^[a-zA-Z].([A-Za-z0-9_ ]){3,50}$`)
var validDescription = regexp.MustCompile(`^[a-zA-Z].([A-Za-z0-9_ \.\?,']){3,300}$`)

// Validate checks the action params for errs
func (p *ActionParams) Validate() error {

	hasErrors := false
	validationErrs := &ActionParams{}

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
func (a *Action) CreateAction(db *sql.DB, params ActionParams, userID string) error {
	stmt, err := db.Prepare("INSERT INTO actions (actionID, title, description, userID) VALUES(UUID_TO_BIN(UUID()), ?, ?, UUID_TO_BIN(?))")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(params.Title, params.Description, userID)
	if err != nil {
		return dbservice.CheckDatabaseErr(err, "title")
	}

	return nil
}

// GetActions retrieves all actions in the db
func (a *Action) GetActions(db *sql.DB) ([]Action, error) {
	stmt, err := db.Prepare("SELECT BIN_TO_UUID(actionID)actionID,title,description,isArchived,createdAt,updatedAt,BIN_TO_UUID(userID)userID FROM actions")
	if err != nil {
		return nil, dbservice.CheckDatabaseErr(err)
	}

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	var actions []Action
	for rows.Next() {
		var action Action

		err := rows.Scan(
			&action.ActionID,
			&action.Title,
			&action.Description,
			&action.isArchived,
			&action.CreatedAt,
			&action.UpdatedAt,
			&action.UserID,
		)

		if err != nil {
			return nil, err
		}
		// ignore archived actions
		if !action.isArchived {
			actions = append(actions, action)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// GetActionByID retrueves a single action by its actionID
func (a *Action) GetActionByID(db *sql.DB, actionID string) (*Action, error) {
	stmt, err := db.Prepare(fmt.Sprintf("SELECT BIN_TO_UUID(actionID)actionID,title,description,isArchived,createdAt,updatedAt,BIN_TO_UUID(userID)userID FROM actions WHERE actionID=UUID_TO_BIN('%v')", actionID))

	if err != nil {
		return nil, err
	}

	action := Action{}
	err = stmt.QueryRow().Scan(
		&action.ActionID,
		&action.Title,
		&action.Description,
		&action.isArchived,
		&action.CreatedAt,
		&action.UpdatedAt,
		&action.UserID,
	)

	if err != nil {
		fmt.Println("ERROR", err)

		return nil, err
	}

	if action.isArchived {
		return nil, sql.ErrNoRows
	}

	return &action, nil
}

// UpdateAction updates an action's title or description
func (a *Action) UpdateAction(db *sql.DB, actionID string, params ActionParams) error {
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
