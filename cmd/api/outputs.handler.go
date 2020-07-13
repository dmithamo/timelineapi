package main

import "net/http"

// createOutput handles requests for creating a new output
// Accessible @ POST /outputs
func (a *application) createOutput(w http.ResponseWriter, r *http.Request) {}

// getOutputs handles requests for retrieving all outputs
// Accessible @ GET /outputs
func (a *application) getOutputs(w http.ResponseWriter, r *http.Request) {}

// getOutput handles requests for retrieving a single output by outputID
// Accessible @ GET /outputs/{outputID}
func (a *application) getOutput(w http.ResponseWriter, r *http.Request) {}

// getOutputsByAction handles requests for retrieving a project's outputs by projectID
// Accessible @ GET /outputs/{projectID}
func (a *application) getOutputsByAction(w http.ResponseWriter, r *http.Request) {}

// updateOutput handles requests for editing output
// Accesible @ PATCH /outputs/{outputID}
func (a *application) updateOutput(w http.ResponseWriter, r *http.Request) {}

// deleteOutput handles requests for deleting a output
// Accessible @ DELETE /outputs/{outputID}
func (a *application) deleteOutput(w http.ResponseWriter, r *http.Request) {}
