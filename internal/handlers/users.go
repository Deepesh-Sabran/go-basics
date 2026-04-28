package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Deepesh-Sabran/go-basics/internal/helper"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
	"github.com/Deepesh-Sabran/go-basics/internal/services"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Invalid request for login")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// get the token by calling the service
	tokens, err:= services.Login(req.Name, req.Password)
	if err != nil {
		log.Println("Unauthorized user")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct{
		RefreshToken string `json:"refreshToken"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Invalid request for generating Refresh Token")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newAccessToken, err:= services.Refresh(req.RefreshToken) 
	if err != nil {
		log.Println("❌ ERROR: Invalid refresh token")
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	response:= map[string]string{
		"accessToken": newAccessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	// convert req to user model because service is using DB model
	user:= models.User{
		Name: 		req.Name,
		Age:		req.Age,
		Password: 	req.Password,
		Role:		req.Role,
	}

	err:= services.CreateUser(&user)
	if err != nil {
		log.Println("❌ ERROR: CreateUser service error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// now what you get from service convert that to UserResponse model
	res:= models.UserResponse{
		ID:		user.ID,
		Name:	user.Name,
		Age:	user.Age,
		Role:	user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	page, limit:= 1, 5

	query:= r.URL.Query()
	// pageStr:= query.Get("page")
	// limitStr:= query.Get("limit")

	if p := query.Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := query.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			if val > 50 {
				limit = 50
			} else {
				limit = val
			}
		}
	}


	userList, total, err:= services.GetUsers(page, limit)
	if err != nil {
		log.Println("❌ ERROR: GetUsers service error: ", err)
		http.Error(w, "Error fetching users", http.StatusBadRequest)
		return
	}

	var users []models.UserResponse

	for _, user:= range userList{
		users = append(users, helper.ToUserResponse(user))
	}

	response:= models.PaginatedUserResponse{
		Data: users,
		Page:  page,
		Limit: limit,
		Count: total,
	}

	userId:= r.Context().Value("user_id")
	name := r.Context().Value("name")
	role:= r.Context().Value("role")

	log.Println("🔥 Request by:", userId, name, role)

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	res:= models.UserResponse{
		ID: 	user.ID,
		Name: 	user.Name,
		Age: 	user.Age,
		Role:	user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userId:= r.Context().Value("user_id")

	idFloat, ok:= userId.(float64)
	if !ok {
		log.Println("❌ ERROR: user id is not valid")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	idInt:= int(idFloat)

	user, err:= services.GetUserById(idInt)
	if err != nil {
		log.Println("❌ User not found")
		http.Error(w, "😞 User not found", http.StatusNotFound)
		return
	}

	response:= helper.ToUserResponse(*user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	// id from token
	tokenIdFloat, okId:= r.Context().Value("user_id").(float64)
	if !okId {
		log.Println("🚫 You are unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
    	return
	}
	tokenIdInt:= int(tokenIdFloat)
	log.Println("tokenIdInt: ", tokenIdInt)

	// check permission from context
	hasUpdatePermission:= false

	perms, ok:= r.Context().Value("permissions").([]interface{})
	if !ok {
		log.Println("you are forbidden to do this action")
		http.Error(w, "User Forbidden", http.StatusForbidden)
		return
	}

	for _, p:= range perms {
		if p.(string) == "update_user" {
			hasUpdatePermission = true
			break
		}
	}

	if !hasUpdatePermission && tokenIdInt != id {
		log.Println("🚫 You are forbidden to do this action")
		http.Error(w, "🚫 You are forbidden to do this action", http.StatusForbidden)
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