package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Deepesh-Sabran/go-basics/internal/models"
	"github.com/Deepesh-Sabran/go-basics/internal/services"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err:= json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	err:= services.CreateUser(&user)
	if err != nil {
		log.Println("❌ ERROR: CreateUser service error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	userList, err:= services.GetUsers()
	if err != nil {
		log.Println("❌ ERROR: GetUsers service error: ", err)
		http.Error(w, "Error fetching users", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&userList)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	idStr:= r.PathValue("id")
	id, err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	user, err:= services.GetUserById(id)
	if err != nil {
		log.Println("❌ ERROR: GetUserById service error: ", err)
		http.Error(w, "Error fetching user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Header", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}

func DeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	err:= services.DeleteAllUsers(r.Context())
	if err != nil {
		log.Println("❌ ERROR: DeleteAllUsers service error: ", err)
		http.Error(w, "Error deleting users", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "All users deleted successfully 🗑️" })
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	idStr:= r.PathValue("id")
	id, err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = services.DeleteUserById(r.Context(), id)
	if err != nil {
		log.Println("❌ ERROR: DeleteUserById service error: ", err)
		http.Error(w, "Error deleting user", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully 🗑️"})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr:= r.PathValue("id")
	id, err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err:= json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
		return
	}

	err = services.UpdateUser(id, updates)
	if err != nil {
		log.Println("❌ ERROR: UpdateUser service error: ", err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Updated Successfully"})
}