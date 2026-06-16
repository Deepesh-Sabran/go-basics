package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	appErrors "github.com/Deepesh-Sabran/go-basics/internal/errors"
	"github.com/Deepesh-Sabran/go-basics/internal/helper"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
	"github.com/Deepesh-Sabran/go-basics/internal/services"
	"github.com/Deepesh-Sabran/go-basics/internal/validation"
	"github.com/Deepesh-Sabran/go-basics/pkg/response"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	
	decoder:= json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err:= decoder.Decode(&req); err != nil {
		log.Println("❌ ERROR: ", err)
		response.HandleError(w, appErrors.BadRequest("invalid request body"))
		return
	}

	if decoder.More() {
		response.HandleError(w, appErrors.BadRequest("invalid request body"))
		return
	}

	if err:= validation.ValidateLogin(req); err != nil {
		response.HandleError(w, err)
		return
	}

	// get the token by calling the service
	tokens, err:= services.Login(req.Name, req.Password)
	if err != nil {
		log.Println("❌ ERROR: ", err)
		response.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest

	decoder:=json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err:= decoder.Decode(&req); err != nil {
		log.Println("Invalid request for generating Refresh Token")
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if decoder.More() {
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if err:= validation.ValidateRefresh(req); err != nil {
		response.HandleError(w, err)
		return
	}

	newAccessToken, err:= services.Refresh(req.RefreshToken) 
	if err != nil {
		log.Println("❌ ERROR: Invalid refresh token")
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]string{
		"access_token": newAccessToken,
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	decoder:= json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err:= decoder.Decode(&req); err != nil {
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if decoder.More() {
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	// validate inputs using validation helper
	if err:= validation.ValidateCreateUser(req); err != nil {
		response.HandleError(w, err)
		return
	}

	// convert req to user model because service is using DB model
	user:= models.User{
		Name: 		req.Name,
		Age:		req.Age,
		Password: 	req.Password,
	}

	err:= services.CreateUser(&user)
	if err != nil {
		log.Println("❌ ERROR: CreateUser service error: ", err)
		response.HandleError(w, err)
		return
	}

	// now what you get from service convert that to UserResponse model
	res:= helper.ToUserResponse(user)

	response.HandleSuccess(w, http.StatusCreated, res)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	page, limit:= 1, 5

	query:= r.URL.Query()

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
		response.HandleError(w, err)
		return
	}

	var users []models.UserResponse

	for _, user:= range userList{
		users = append(users, helper.ToUserResponse(user))
	}

	respData:= models.PaginatedUserResponse{
		Data: users,
		Page:  page,
		Limit: limit,
		Count: total,
	}

	userId:= r.Context().Value("user_id")
	name := r.Context().Value("name")

	log.Printf("🔥 Request by==> id: %d  name: %s", userId, name)

	response.HandleSuccess(w, http.StatusOK, respData)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	idStr:= r.PathValue("id")
	id, err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID format")
		response.HandleError(w, appErrors.BadRequest("Invalid ID format"))
		return
	}

	user, err:= services.GetUserById(id)
	if err != nil {
		log.Println("❌ ERROR: ", err)
		response.HandleError(w, err)
		return
	}

	res:= models.UserResponse{
		ID: 	user.ID,
		Name: 	user.Name,
		Age: 	user.Age,
		RoleID:	user.RoleID,
	}

	response.HandleSuccess(w, http.StatusOK, res)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userId:= r.Context().Value("user_id").(int)

	idInt:= int(userId)

	user, err:= services.GetUserById(idInt)
	if err != nil {
		log.Println("❌ User not found")
		response.HandleError(w, err)
		return
	}

	respUser:= helper.ToUserResponse(*user)

	response.HandleSuccess(w, http.StatusOK, respUser)
}

func DeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	err:= services.DeleteAllUsers(r.Context())
	if err != nil {
		log.Println("❌ ERROR: DeleteAllUsers service error: ", err)
		response.HandleError(w, appErrors.InternalServerError("Error deleting users"))
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]string{"message": "All users deleted successfully 🗑️" })
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	// take ID from path param
	idStr:= r.PathValue("id")
	pathId, err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID format")
		response.HandleError(w, appErrors.BadRequest("Invalid ID format"))
		return
	}

	// take ID from context
	tokenId, ok:= r.Context().Value("user_id").(int)
	if !ok {
		response.HandleError(w, appErrors.Unauthorized("Unauthorized"))
		return
	}

	// call new DeleteUserWithAudit function
	if err:= services.DeleteUserWithAudit(r.Context(), tokenId, pathId); err != nil {
		log.Println("❌ ERROR: DeleteUserWithAudit service error: ", err)
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]string{"message": "User deleted successfully 🗑️"})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateUserRequest

	idStr:= r.PathValue("id")
	id, err:= strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID format")
		response.HandleError(w, appErrors.BadRequest("Invalid ID format"))
		return
	}

	decoder:= json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err:= decoder.Decode(&req); err != nil {
		log.Println("update decode error:", err)
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if decoder.More() {
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if err:= validation.ValidateUpdateUser(req); err != nil {
		response.HandleError(w, err)
		return
	}

	updates := req.ToUpdatesMap()

	err = services.UpdateUser(id, updates)
	if err != nil {
		log.Println("❌ ERROR: UpdateUser service error: ", err)
		response.HandleError(w, err)
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]string{"message": "User Updated Successfully"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest

	decoder:= json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err:= decoder.Decode(&req); err != nil {
		log.Println("Invalid request for logout")
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if decoder.More() {
		response.HandleError(w, appErrors.BadRequest("Invalid request body"))
		return
	}

	if err:= validation.ValidateRefresh(req); err != nil {
		response.HandleError(w, err)
		return
	}

	if err:= services.Logout(req.RefreshToken); err != nil {
		log.Println("Error during logout")
		response.HandleError(w, err)
		return
	}

	log.Println("Logged out successfully")

	response.HandleSuccess(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}