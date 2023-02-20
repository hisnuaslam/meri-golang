package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// User is a struct that represents a user in the database
type User struct {
	DoctorID         int    `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
	DoctorSpecialist string `json:"specialist"`
}

func main() {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "sameri:S4meri@p2!@tcp(172.16.1.8:3306)/PRI_TX_MERI_LAB")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Define the routes for the API
	r.HandleFunc("/users", getUsers(db)).Methods("GET")
	r.HandleFunc("/users/{doctor_id}", getUser(db)).Methods("GET")
	r.HandleFunc("/users", createUser(db)).Methods("POST")
	r.HandleFunc("/users/{doctor_id}", updateUser(db)).Methods("PUT")
	r.HandleFunc("/users/{doctor_id}", deleteUser(db)).Methods("DELETE")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// getUsers is a handler function that returns a list of all users in the database
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database for all users
		rows, err := db.Query("SELECT doctor_id, doctor_name, specialist FROM MS_DOCTOR")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate through the rows and build a slice of users
		users := []User{}
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.DoctorID, &user.DoctorName, &user.DoctorSpecialist); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		// Convert the slice of users to JSON and write it to the response
		json.NewEncoder(w).Encode(users)
	}
}

// getUser is a handler function that returns a specific user by ID
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the request URL
		vars := mux.Vars(r)
		fmt.Println(vars)
		id := vars["doctor_id"]
		fmt.Println(id)

		// Query the database for the user with the specified ID
		row := db.QueryRow("SELECT doctor_id, doctor_name, specialist FROM MS_DOCTOR WHERE doctor_id = ?", id)
		// Build a user struct from the row and write it to the response
		var user User
		if err := row.Scan(&user.DoctorID, &user.DoctorName, &user.DoctorSpecialist); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

// createUser is a handler function that creates a new user in the database
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the JSON request body into a user struct
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insert the user into the database and get the ID of the new row
		_, err := db.Exec("INSERT INTO MS_DOCTOR (doctor_id, doctor_name, specialist) VALUES (?, ?, ?)", user.DoctorID, user.DoctorName, user.DoctorSpecialist)
		// result, err := db.Exec("INSERT INTO MS_DOCTOR (doctor_name, specialist) VALUES (?, ?)", user.DoctorName, user.DoctorSpecialist)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// id, _ := result.LastInsertId()

		// Set the ID of the user struct and write it to the response
		// user.DoctorID = int(id)
		json.NewEncoder(w).Encode(user)
	}
}

// updateUser is a handler function that updates a user in the database
// updateUser is a handler function that updates a user in the database
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the request URL
		vars := mux.Vars(r)
		fmt.Println(vars)
		id, err := strconv.Atoi(vars["doctor_id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(id)
		// Decode the JSON request body into a user struct
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// // Update the user in the database
		// // user.DoctorID = int(id)

		_, err = db.Exec("UPDATE MS_DOCTOR SET doctor_name = ?, specialist = ? WHERE doctor_id = ?", user.DoctorName, user.DoctorSpecialist, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// // Set the ID of the user struct and write it to the response
		user.DoctorID = int(id)
		json.NewEncoder(w).Encode(user)
	}
}

// // deleteUser is a handler function that deletes a user from the database
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the request URL
		vars := mux.Vars(r)
		id := vars["doctor_id"]
		fmt.Println(vars)

		// Delete the user from the database
		_, err := db.Exec("DELETE FROM MS_DOCTOR WHERE doctor_id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write a success message to the response
		message := map[string]string{"message": "User deleted successfully"}
		json.NewEncoder(w).Encode(message)
	}
}
