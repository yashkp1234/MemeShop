package controller

import "net/http"

// GetUsers lists all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List Users"))
}

// GetUser lists a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List User"))
}

// CreateUser creates a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create User"))
}

// UpdateUser updates a user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update User"))
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete User"))
}
