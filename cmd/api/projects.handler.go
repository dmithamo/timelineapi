package main

import "net/http"

// createProject handles requests for creating a new project
// Accessible @ POST /projects
func (a *application) createProject(w http.ResponseWriter, r *http.Request) {}

// getProjects handles requests for retrieving all projects
// Accessible @ GET /projects
func (a *application) getProjects(w http.ResponseWriter, r *http.Request) {}

// getProject handles requests for retrieving a single project by projectID
// Accessible @ GET /projects/{projectID}
func (a *application) getProject(w http.ResponseWriter, r *http.Request) {}

// updateProject handles requests for editing project
// Accesible @ PATCH /projects/{projectID}
func (a *application) updateProject(w http.ResponseWriter, r *http.Request) {}

// deleteProject handles requests for deleting a project
// Accessible @ DELETE /projects/{projectID}
func (a *application) deleteProject(w http.ResponseWriter, r *http.Request) {}
