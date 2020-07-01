package main

import (
	"net/http"
)

// register handles requests for creating new users
// Accessible @ POST /auth/register
func registerUser(w http.ResponseWriter, r *http.Request) {}

// updateUser handles request for editing user
// Accessible @ PATCH /auth/register
func updateUser(w http.ResponseWriter, r *http.Request) {}

// login handles requests for login
// Accessible @ POST /auth/login
func loginUser(w http.ResponseWriter, r *http.Request) {}
