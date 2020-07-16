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
	credentials, ok := a.decodeParamsHelper(w, r)

	if !ok {
		// err is already handled (sent back as jsonRes to user)
		return
	}

	var u users.Model
	err := u.CreateUser(a.db, credentials)

	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, &utils.GenericJSONRes{
			Message: "err creating user",
			Data:    AuthError{err.Error()},
		})

		return
	}

	// success!
	// TODO: login user after successful registration?
	utils.SendJSONResponse(w, http.StatusCreated, &utils.GenericJSONRes{
		Message: "successfully created user",
		Data:    nil,
	})
}

// login handles requests for login
// Accessible @ POST /auth/login
func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	credentials, ok := a.decodeParamsHelper(w, r)

	if !ok {
		// err is already handled (sent back as jsonRes to user)
		return
	}

	var u users.Model
	user, err := u.GetByCredentials(a.db, credentials)

	if err != nil {
		if err == sql.ErrNoRows || err.Error() == utils.WRONG_PASSWORD_ERR {
			utils.SendJSONResponse(w, http.StatusBadRequest, &utils.GenericJSONRes{
				Message: "wrong username or password",
				Data:    nil,
			})

			return
		}

		// any other errs
		utils.SendJSONResponse(w, http.StatusInternalServerError, &utils.GenericJSONRes{
			Message: fmt.Sprintf("err loggin in: %v", err.Error()),
			Data:    nil,
		})

		return
	}

	// if no err, sign user in
	a.loginUserHelper(w, r, user)
}

// decodeParamsHelper decodes request body into a credentials struct
func (a *application) decodeParamsHelper(w http.ResponseWriter, r *http.Request) (*users.Credentials, bool) {
	var credentials users.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)

	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, &utils.GenericJSONRes{
			Message: fmt.Sprintf("err decoding request body: %v", err.Error()),
			Data:    nil,
		})

		return nil, false
	}

	validationErrs := credentials.Validate()
	if validationErrs != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, &utils.GenericJSONRes{
			Message: "invalid user credentials",
			Data:    validationErrs,
		})

		return nil, false
	}

	return &credentials, true
}

// loginUserHelper logs a user in
func (a *application) loginUserHelper(w http.ResponseWriter, r *http.Request, user *users.Model) {
	// generate session token, create a cookie for the client
	token := uuid.New().String()

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(60 * time.Second),
		HttpOnly: true,
		Path:     "/",
		SameSite: 1,
	})

	// success!
	utils.SendJSONResponse(w, http.StatusOK, &utils.GenericJSONRes{
		Message: "successfully logged in",
		Data:    nil,
	})
}

// updateUser handles request for editing user
// Accessible @ PATCH /auth/register
func (a *application) updateUser(w http.ResponseWriter, r *http.Request) {}
