package handlers

import (
	"encoding/json"
	"net/http"
)

type User struct{
	Name string `json:"name"`
	Age  int	`json:"age"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	// decode request body so that your GO understand the data
	if err:= json.NewDecoder(r.Body).Decode(&user); err != nil {		// once you read the "r.Body" next time it is blank
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// check if any payload value is missing then through error
	if user.Name == "" || user.Age == 0 {
		http.Error(w, "Wrong Payload", http.StatusBadRequest)
		return
	}

	// create response
	response:= map[string]interface{}{
		"message": "User created successfully 🚀",
		"data"	 : user,
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode response to JSON so that browser can understand
	json.NewEncoder(w).Encode(response)
}