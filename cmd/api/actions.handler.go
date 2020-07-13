package main

import "net/http"

// createAction handles requests for creating a new action
// Accessible @ POST /actions
func (a *application) createAction(w http.ResponseWriter, r *http.Request) {}

// getActions handles requests for retrieving all Actions
// Accessible @ GET /actions
func (a *application) getActions(w http.ResponseWriter, r *http.Request) {}

// getAction handles requests for retrieving a single Action by ActionID
// Accessible @ GET /actions/{actionID}
func (a *application) getAction(w http.ResponseWriter, r *http.Request) {}

// updateAction handles requests for editing Action
// Accesible @ PATCH /actions/{actionID}
func (a *application) updateAction(w http.ResponseWriter, r *http.Request) {}

// deleteAction handles requests for deleting a Action
// Accessible @ DELETE /actions/{actionID}
func (a *application) deleteAction(w http.ResponseWriter, r *http.Request) {}
