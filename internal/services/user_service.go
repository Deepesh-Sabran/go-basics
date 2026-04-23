package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Deepesh-Sabran/go-basics/internal/models"
	repo "github.com/Deepesh-Sabran/go-basics/internal/repository"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("VENGEANCE") // letter move to environment file

func Login(name, password string) (map[string]string, error) {
	// call GetUserByName repo function to get the user
	user, err:= repo.GetUserByName(name)
	if err != nil {
		log.Println("User not found 😞")
		return nil, errors.New("User not found")
	}

	// check & compare password entered by user and stored in DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid password")
		return nil, errors.New("Incorrect password, please check your password")
	}

	// create access token
	accessClaims:= jwt.MapClaims{
		"user_id":		user.ID,
		"name":			user.Name,
		"exp":			time.Now().Add(time.Minute * 15).Unix(),
	}

	accessToken:= jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err:= accessToken.SignedString(jwtSecret)
	if err != nil {
		log.Println("❌ ERROR: Invalid access token", err)
        return nil, err
    }

	// create refresh token
	refreshClaims:= jwt.MapClaims{
		"user_id":		user.ID,
		"exp":			time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken:= jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err:= refreshToken.SignedString(jwtSecret)
	if err != nil {
		log.Println("❌ ERROR: Invalid refresh token", err)
        return nil, err
    }

	log.Println("User LogIn successful 🥳")

	return map[string]string{
		"access token": accessString,
		"refresh token": refreshString,
	}, nil
}

func Refresh(refreshToken string) (string, error) {
	token, err:= jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		log.Println("❌ ERROR: Invalid refresh token")
		return "", errors.New("Invalid refresh token, login again")
	}

	claims, ok:= token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid Claims")
	}

	userId:= claims["user_id"]
	
	// create new access token
	newClaims:= jwt.MapClaims{
		"user_id":	userId,
		"exp"	 :	time.Now().Add(time.Minute * 15).Unix(),
	}

	newToken:= jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err:= newToken.SignedString(jwtSecret)
	if err != nil {
		log.Println("❌ ERROR: Invalid access token", err)
        return "", err
    }

	return newTokenString, nil
}

func CreateUser(user *models.User) error {
	if user.Name == "" || user.Age <= 0 || user.Password == "" {
		return errors.New("Invalid user data")
	}

	hashedPassword, err:= bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ ERROR: failed to hash password: ", err)
		return err
	}

	user.Password = string(hashedPassword)

	err = repo.CreateUser(user)
	if err != nil {
		log.Println("💔 Failed to create user: ", err)
		return err
	}

	log.Println("User SignUp successful 🥳")
	return nil
}

func GetUsers(page, limit int) ([]models.User, error) {
	offset:= (page - 1) * limit

	userList, err:= repo.GetUsers(limit, offset)
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

	if val, ok := updates["password"]; ok {
		passwordStr, ok := val.(string)
		if !ok {
			return errors.New("Invalid password format")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ ERROR: failed to hash password:", err)
			return err
		}

		cleanUpdates["password"] = string(hashedPassword) // store as string
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