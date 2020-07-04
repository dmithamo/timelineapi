package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmithamo/timelineapi/pkg/users"
	"github.com/dmithamo/timelineapi/pkg/utils"
)

// AuthError defines a generic auth error
type AuthError struct {
	Message string `json:"detail,omitempty"`
}

// AuthorizationToken defines the data returned after successful auth event
type AuthorizationToken struct {
	Token string      `json:"token,omitempty"`
	User  users.Model `json:"user,omitempty"`
}

// register handles requests for creating new users
// Accessible @ POST /auth/register
func (a *application) registerUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var credentials users.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err decoding request body: %v", err.Error()), nil)
		return
	}

	validationErrs := credentials.Validate()
	if validationErrs != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "invalid user credentials", validationErrs)
		return
	}

	var u users.Model
	err = u.CreateUser(a.db, &credentials)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "unable to create user", AuthError{err.Error()})
		return
	}

	// success!
	utils.SendJSONResponse(w, http.StatusCreated, "successfully created user", nil)
}

// login handles requests for login
// Accessible @ POST /auth/login
func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// for security reasons, exact err message is dropped in
	// favor of obfuscatedErrMessage
	obfuscatedErrMessage := AuthError{"wrong username or password"}

	var credentials users.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err decoding request body: %v", err.Error()), nil)
		return
	}

	validationErrs := credentials.Validate()
	if validationErrs != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "invalid user credentials", obfuscatedErrMessage)
		return
	}

	// proceed with login
	var u users.Model
	user, err := u.GetByCredentials(a.db, &credentials)
	if err != nil {
		if err == sql.ErrNoRows {
			// credentials no match!
			utils.SendJSONResponse(w, http.StatusBadRequest, "invalid user credentials", obfuscatedErrMessage)
			return
		}
		// any other errs
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err loging in: %v", err.Error()), nil)
		return
	}

	// success!
	utils.SendJSONResponse(w, http.StatusOK, "successfully logged in", generateAuthToken(user))
}

// generateAuthToken generates an auth token and returns it with user embedded
func generateAuthToken(user *users.Model) interface{} {
	return AuthorizationToken{"fake-token-for-now", *user}
}

// updateUser handles request for editing user
// Accessible @ PATCH /auth/register
func (a *application) updateUser(w http.ResponseWriter, r *http.Request) {}
