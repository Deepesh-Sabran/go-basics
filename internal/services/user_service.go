package services

import (
	"context"
	"errors"
	"log"

	"github.com/Deepesh-Sabran/go-basics/internal/models"
	repo "github.com/Deepesh-Sabran/go-basics/internal/repository"
)

func CreateUser(user *models.User) error {
	if user.Name == "" || user.Age <= 0 {
		return errors.New("Invalid user data")
	}

	err:= repo.CreateUser(user)
	if err != nil {
		log.Println("💔 Failed to create user: ", err)
		return err
	}

	log.Println("User Created successfully 🥳")
	return nil
}

func GetUsers() ([]models.User, error) {
	userList, err:= repo.GetUsers()
	if err != nil {
		log.Println("💔 Failed to fetch user list")
		return nil, err
	}

	log.Println("🥳 User List fetched successfully")
	return userList, nil
}

func GetUserById(userId int) (*models.User, error) {
	user, err:= repo.GetUserById(userId)
	if err != nil {
		log.Println("💔 Failed to fetch a user")
		return nil, err
	}

	if user == nil {
		log.Println("😞 User not found")
		return nil, errors.New("User not found")
	}

	log.Println("🥳 User fetched successfully")
	return user, nil
}

func DeleteAllUsers(ctx context.Context) error {
	err:= repo.DeleteAllUsers(ctx)
	if err != nil {
		log.Println("🚫 Failed to delete users")
		return err
	}

	log.Println("😉 Users deleted successfully")
	return nil
}

func DeleteUserById(ctx context.Context, userId int) error {
	err:= repo.DeleteUserById(ctx, userId)
	if err != nil {
		log.Println("🚫 Failed to delete the user")
		return err
	}

	log.Println("😉 User deleted successfully")
	return nil
}

func UpdateUser(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("No data provided")
	}

	// white listing the fields which are allowed
	cleanUpdates:= make(map[string]interface{})

	if val, ok:= updates["name"]; ok {
		cleanUpdates["name"] = val
	}

	if val, ok:= updates["age"]; ok {
		cleanUpdates["age"] = val
	}

	// after cleaning if it's empty ...
	if len(cleanUpdates) == 0 {
		return errors.New("Invalid fields provided")
	}

	err:= repo.UpdateUser(id, cleanUpdates)
	if err != nil {
		log.Println("❌ ERROR: In updating users")
		return err
	}

	log.Println("🥳 User updated successfully")
	return nil
}