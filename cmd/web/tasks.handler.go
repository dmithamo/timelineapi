package main

import "net/http"

// createTask handles requests for creating a new task
// Accessible @ POST /tasks
func createTask(w http.ResponseWriter, r *http.Request) {}

// getTasks handles requests for retrieving all Tasks
// Accessible @ GET /tasks
func getTasks(w http.ResponseWriter, r *http.Request) {}

// getTask handles requests for retrieving a single task by taskID
// Accessible @ GET /tasks/{taskID}
func getTask(w http.ResponseWriter, r *http.Request) {}

// getTasksByProject handles requests for retrieving a project's tasks by projectID
// Accessible @ GET /tasks/{projectID}
func getTasksByProject(w http.ResponseWriter, r *http.Request) {}

// updateTask handles requests for editing Task
// Accesible @ PATCH /tasks/{taskID}
func updateTask(w http.ResponseWriter, r *http.Request) {}

// deleteTask handles requests for deleting a task
// Accessible @ DELETE /tasks/{taskID}
func deleteTask(w http.ResponseWriter, r *http.Request) {}
