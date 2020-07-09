package main

import "net/http"

// createTask handles requests for creating a new task
// Accessible @ POST /tasks
func (a *application) createTask(w http.ResponseWriter, r *http.Request) {}

// getTasks handles requests for retrieving all Tasks
// Accessible @ GET /tasks
func (a *application) getTasks(w http.ResponseWriter, r *http.Request) {}

// getTask handles requests for retrieving a single task by taskID
// Accessible @ GET /tasks/{taskID}
func (a *application) getTask(w http.ResponseWriter, r *http.Request) {}

// getTasksByProject handles requests for retrieving a project's tasks by projectID
// Accessible @ GET /tasks/{projectID}
func (a *application) getTasksByProject(w http.ResponseWriter, r *http.Request) {}

// updateTask handles requests for editing Task
// Accesible @ PATCH /tasks/{taskID}
func (a *application) updateTask(w http.ResponseWriter, r *http.Request) {}

// deleteTask handles requests for deleting a task
// Accessible @ DELETE /tasks/{taskID}
func (a *application) deleteTask(w http.ResponseWriter, r *http.Request) {}
