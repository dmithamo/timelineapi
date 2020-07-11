package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dmithamo/timelineapi/pkg/users"
	"github.com/dmithamo/timelineapi/pkg/utils"
	"github.com/google/uuid"
)

// AuthError defines a generic auth error
type AuthError struct {
	Message string `json:"detail,omitempty"`
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
	// TODO: login user after successful registration?
	utils.SendJSONResponse(w, http.StatusCreated, "successfully created user", nil)
}

// login handles requests for login
// Accessible @ POST /auth/login
func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	loginErrMessage := AuthError{"wrong username or password"}

	var credentials users.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err decoding request body: %v", err.Error()), nil)
		return
	}

	validationErrs := credentials.Validate()
	if validationErrs != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "invalid user credentials", loginErrMessage)
		return
	}

	var u users.Model
	user, err := u.GetByCredentials(a.db, &credentials)
	if err != nil {
		if err == sql.ErrNoRows || err.Error() == utils.WRONG_PASSWORD_ERR {
			// credentials no match!
			utils.SendJSONResponse(w, http.StatusBadRequest, "invalid user credentials", loginErrMessage)
			return
		}
		// any other errs
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err loging in: %v", err.Error()), nil)
		return
	}

	// if no err, sign user in
	a.loginUserHelper(w, r, user)
}

// loginUserHelper logs a user in
func (a *application) loginUserHelper(w http.ResponseWriter, r *http.Request, user *users.Model) {
	// SESSION_TOKEN_LIFETIME is the time in seconds until expiration of a session_token
	const SESSION_TOKEN_LIFETIME = 60 * 60

	// generate session token, create a cookie for the client
	token := uuid.New().String()
	_, err := a.cache.Do("SETEX", token, SESSION_TOKEN_LIFETIME, user.UserID)

	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err creating session: %v", err.Error()), nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(SESSION_TOKEN_LIFETIME * time.Second),
		HttpOnly: true,
		Path:     "/",
		SameSite: 1,
	})

	// success!
	utils.SendJSONResponse(w, http.StatusOK, "successfully logged in", nil)
}

// updateUser handles request for editing user
// Accessible @ PATCH /auth/register
func (a *application) updateUser(w http.ResponseWriter, r *http.Request) {}
