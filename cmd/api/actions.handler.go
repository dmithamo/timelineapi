package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmithamo/timelineapi/pkg/actions"
	"github.com/dmithamo/timelineapi/pkg/utils"
	"github.com/gorilla/mux"
)

// ActionErr structures an err that arises during action
type ActionErr struct {
	Message string `json:"detail,omitempty"`
}

// createAction handles requests for creating a new action
// Accessible @ POST /actions
func (a *application) createAction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var actionParams actions.Params
	err := json.NewDecoder(r.Body).Decode(&actionParams)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("err decoding request body: %v", err.Error()), nil)
		return
	}

	validationErrs := actionParams.Validate()
	if validationErrs != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "invalid action params", validationErrs)
		return
	}

	var actionModel actions.Model
	err = actionModel.CreateAction(a.db, actionParams)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "unable to create action", ActionErr{err.Error()})
		return
	}

	// success!
	utils.SendJSONResponse(w, http.StatusCreated, "successfully created action", nil)
}

// getActions handles requests for retrieving all Actions
// Accessible @ GET /actions
func (a *application) getActions(w http.ResponseWriter, r *http.Request) {
	var actionModel actions.Model

	actions, err := actionModel.GetActions(a.db)
	if err != nil || actions == nil {
		if err == sql.ErrNoRows || actions == nil {
			// TODO: no actions found. Should this be a 404?
			utils.SendJSONResponse(w, http.StatusNotFound, "no actions found", nil)
			return
		}

		utils.SendJSONResponse(w, http.StatusBadRequest, "unable to retrieve actions", ActionErr{err.Error()})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "successfully retrieved actions", actions)
}

// getAction handles requests for retrieving a single Action by ActionID
// Accessible @ GET /actions/{actionID}
func (a *application) getAction(w http.ResponseWriter, r *http.Request) {
	var actionModel actions.Model

	actionID := mux.Vars(r)["actionID"]

	action, err := actionModel.GetActionByID(a.db, actionID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("no action found with actionID `%v`", actionID), nil)
			return
		}

		utils.SendJSONResponse(w, http.StatusBadRequest, "unable to retrieve action", ActionErr{err.Error()})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "successfully retrieved actions", action)
}

// updateAction handles requests for editing Action
// Accesible @ PATCH /actions/{actionID}
func (a *application) updateAction(w http.ResponseWriter, r *http.Request) {}

// deleteAction handles requests for deleting a Action
// Accessible @ DELETE /actions/{actionID}
func (a *application) deleteAction(w http.ResponseWriter, r *http.Request) {}
