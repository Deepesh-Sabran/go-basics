package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	appErrors "github.com/Deepesh-Sabran/go-basics/internal/errors"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
	repo "github.com/Deepesh-Sabran/go-basics/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	// "github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)

func Login(name, password string) (map[string]string, error) {
	// call GetUserByName repo function to get the user
	user, err:= repo.GetUserByName(name)
	if err != nil {
		log.Println("User not found 😞")
		return nil, appErrors.NotFound("User not found")
	}

	// check & compare password entered by user and stored in DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid credentials")
		return nil, appErrors.Unauthorized("invalid credentials")
	}

	accessClaims:= models.TokenClaims{
		UserId:			user.ID,
		Name:			user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(config.GetJWTSecret())
	if err != nil {
		log.Println("❌ ERROR: ", err)
		return nil, appErrors.InternalServerError("failed to generate access token")
	}

	// create refresh token
	refreshClaims:= models.TokenClaims{
		UserId: 		user.ID,
		Name:		user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(config.GetJWTSecret())
	if err != nil {
		log.Println("❌ ERROR: ", err)
		return nil, appErrors.InternalServerError("failed to generate refresh token")
	}

	cacheRefreshToken:= fmt.Sprintf("token:refresh:%s", refreshString)

	config.RedisClient.Set(
		config.Ctx,
		cacheRefreshToken,
		user.ID,
		7*24*time.Hour,
	)

	log.Println("User LogIn successful 🥳")

	return map[string]string{
		"access_token": accessString,
		"refresh_token": refreshString,
	}, nil
}

func Refresh(refreshToken string) (string, error) {
	// check refresh token is there in Redis or not
	cacheRefreshToken:= fmt.Sprintf("token:refresh:%s", refreshToken)
	_, err:= config.RedisClient.Get(config.Ctx, cacheRefreshToken).Result()
	if err == redis.Nil {
		log.Println("refresh token not found in redis")
		return "", errors.New("invalid refresh token, login again")
	}
	if err != nil {
		log.Println("redis error:", err)
		return "", err
	}

	// parsing token with claims
	token, err := jwt.ParseWithClaims(refreshToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return config.GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		config.RedisClient.Del(
			config.Ctx,
			cacheRefreshToken,
		)
		log.Println("❌ ERROR: Invalid refresh token")
		return "", errors.New("Invalid refresh token, login again")
	}

	claims, ok:= token.Claims.(*models.TokenClaims)
	if !ok {
		config.RedisClient.Del(config.Ctx, cacheRefreshToken)
		return "", errors.New("Invalid Claims")
	}

	// fetch fresh user
	user, err:= repo.GetUserById(claims.UserId)
	if err != nil {
		log.Println("💔 Failed to fetch a user")
		return "", err
	}

	newClaims:= models.TokenClaims{
		UserId:			user.ID,
		Name:			user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(config.GetJWTSecret())
	if err != nil {
		log.Println("❌ ERROR: Invalid access token", err)
		return "", err
	}

	log.Println("New Access token generated 🥳")
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
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return errors.New("Username already taken")
			}
    	}
		return err
	}

	// background job using goroutine
	// go func(userId int, name string) {
	// 	log.Printf("Background job: user created id: %d and name: %s", userId, name)
	// }(user.ID, user.Name)

	// background job using redis
	job:= models.EmailJob{
		UserID: user.ID,
		Name: 	user.Name,
		Email: 	"demo@example.com",
	}

	// marshal the data --> Producer
	data, err:= json.Marshal(job)
	if err == nil {
		config.RedisClient.LPush(
			config.Ctx,
			"email_jobs",
			data,
		)
	}

	log.Println("User SignUp successful 🥳")
	return nil
}

func GetPermissionsByRole(roleId int) ([]string, error) {
	cacheKey:= fmt.Sprintf("permissions:role:%d", roleId)

	// check redis
	cached, err:= config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var perms []string
		if err:= json.Unmarshal([]byte(cached), &perms); err == nil {
			return perms, nil
		}
	}

	perms, err:= repo.GetPermissionsByRole(roleId)
	if err != nil {
		log.Println("💔 Failed to fetch permissions")
		return nil, err
	}

	// stores in Redis
	data, err:= json.Marshal(perms)
	if err != nil {
		log.Println("Error in Marshaling permissions during storing it into Redis: ", err)
		return nil, err
	}

	config.RedisClient.Set(
		config.Ctx,
		cacheKey,
		data,
		10*time.Minute,
	)

	return perms, nil
}

func GetUsers(page, limit int) ([]models.User, int,  error) {
	offset:= (page - 1) * limit

	userList, err:= repo.GetUsers(limit, offset)
	if err != nil {
		log.Println("💔 Failed to fetch user list")
		return nil, 0, err
	}

	userCount, err:= repo.GetUserCount()
	if err != nil {
		log.Println("💔 Failed to get user count")
		return nil, 0, err
	}

	log.Println("🥳 User List fetched successfully")
	return userList, userCount, nil
}

func GetUserById(userId int) (*models.User, error) {
	user, err:= repo.GetUserById(userId)
	if err != nil {
		log.Println("💔 Failed to fetch a user")
		return nil, appErrors.NotFound("failed to fetch user")
	}

	if user == nil {
		log.Println("😞 User not found")
		return nil, appErrors.NotFound("User not found")
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

func DeleteUserWithAudit(ctx context.Context, actorId, targetId int) error {
	// begin the transaction
	tx, err:= config.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error occurred, on beginning of transaction")
		return err
	}

	defer tx.Rollback()

	// call delete transaction repo function
	if err:= repo.DeleteUserTx(ctx, tx, targetId); err != nil {
		log.Println("Error deleting user within transaction")
		return err
	}

	// call create audit transaction repo function
	if err := repo.CreateAuditLogTx(ctx, tx, actorId, "delete_user", targetId); err != nil {
		log.Println("Error creating audit for user deletion")
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("Transaction Incomplete")
		return err
	}

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

	if err:= repo.UpdateUser(id, cleanUpdates); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return errors.New("username already taken")
			}
		}

		return err
	}

	log.Println("🥳 User updated successfully")
	return nil
}

func Logout(refreshToken string) error {
	cacheRefreshToken:= fmt.Sprintf("token:refresh:%s", refreshToken)

	if err:= config.RedisClient.Del(config.Ctx, cacheRefreshToken).Err(); err != nil {
		log.Println("redis error:", err)
		return err
	}

	// Removed this so that the logout is been idempotent with res:200 OK
	// if deleted == 0 {
	// 	log.Println("Refresh token not found")
	// 	return errors.New("refresh token not found")
	// }

	log.Println("Logged out successfully")
	return nil
}