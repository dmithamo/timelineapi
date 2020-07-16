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

	decodeErr := json.NewDecoder(r.Body).Decode(&actionParams)
	if decodeErr != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, &utils.GenericJSONRes{
			Message: fmt.Sprintf("err decoding request body: %v", decodeErr.Error()),
			Data:    nil,
		})
		return
	}

	validationErrs := actionParams.Validate()
	if validationErrs != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, &utils.GenericJSONRes{
			Message: "validation errors in params",
			Data:    validationErrs,
		})
		return
	}

	var m actions.Model

	createActionErr := m.CreateAction(a.db, actionParams)
	if createActionErr != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, &utils.GenericJSONRes{
			Message: "err creating action",
			Data:    ActionErr{createActionErr.Error()},
		})
		return
	}

	// success!
	utils.SendJSONResponse(w, http.StatusCreated, &utils.GenericJSONRes{
		Message: "successfully created action",
		Data:    nil,
	})
}

// getActions handles requests for retrieving all Actions
// Accessible @ GET /actions
func (a *application) getActions(w http.ResponseWriter, r *http.Request) {
	var actionModel actions.Model

	allActions, err := actionModel.GetActions(a.db)
	if err != nil || allActions == nil {
		if err == sql.ErrNoRows || allActions == nil {
			utils.SendJSONResponse(w, http.StatusNotFound, &utils.GenericJSONRes{
				Message: "no actions found",
				Data:    nil,
			})
			return
		}

		utils.SendJSONResponse(w, http.StatusInternalServerError, &utils.GenericJSONRes{
			Message: "err retrieving actions",
			Data:    ActionErr{err.Error()},
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, &utils.GenericJSONRes{
		Message: "successfully retrieved actions",
		Data:    allActions,
	})
}

// getAction handles requests for retrieving a single Action by ActionID
// Accessible @ GET /actions/{actionID}
func (a *application) getAction(w http.ResponseWriter, r *http.Request) {
	var actionModel actions.Model

	actionID := mux.Vars(r)["actionID"]
	action, err := actionModel.GetActionByID(a.db, actionID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendJSONResponse(w, http.StatusNotFound,
				&utils.GenericJSONRes{
					Message: fmt.Sprintf("no actions found with actionID: %v", actionID),
					Data:    nil,
				})

			return
		}

		utils.SendJSONResponse(w, http.StatusInternalServerError, &utils.GenericJSONRes{
			Message: "err retrieving action",
			Data:    ActionErr{err.Error()},
		})

		return
	}

	utils.SendJSONResponse(w, http.StatusOK, &utils.GenericJSONRes{
		Message: "successfully retrieved action",
		Data:    action,
	})
}

// updateAction handles requests for editing Action
// Accesible @ PATCH /actions/{actionID}
func (a *application) updateAction(w http.ResponseWriter, r *http.Request) {}

// deleteAction handles requests for deleting a Action
// Accessible @ DELETE /actions/{actionID}
func (a *application) deleteAction(w http.ResponseWriter, r *http.Request) {}
