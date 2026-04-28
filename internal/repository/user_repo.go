package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
)

func CreateUser(user *models.User) error {	
	query:= "INSERT INTO users(name, age, password) VALUES($1, $2, $3) RETURNING id"
	return config.DB.QueryRow(query, user.Name, user.Age, user.Password).Scan(&user.ID)
}

func GetUsers(limit, offset int) ([]models.User, error) {
	rows, err:= config.DB.Query("SELECT id, name, age, role FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Println("Get all users query failed to execute")
		return nil, err
	}
	defer rows.Close()

	var userList []models.User

	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Age, &u.Role)
		userList = append(userList, u)
	}

	return userList, nil
}

func GetUserById(userId int) (*models.User, error) {
	var u models.User

	err:= config.DB.QueryRow("SELECT id, name, age, role FROM users WHERE id = $1", userId).Scan(&u.ID, &u.Name, &u.Age, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with ID %d not found", userId)
			return nil, nil
		}
		log.Println("Database error:", err)
		return nil, err
	}

	return &u, nil
}

func GetUserByName(name string) (*models.User, error) {
	query:= "SELECT id, name, age, password, role FROM users WHERE name=$1"

	var u models.User
	err:= config.DB.QueryRow(query, name).Scan(&u.ID, &u.Name, &u.Age, &u.Password, &u.Role)
	if err != nil {
		log.Printf("User with Name: %s is not fond", name)
		return nil, err
	}

	log.Println("User found with Name: ", name)
	return &u, nil
}

func GetUserCount() (int, error) {
	var count int
	query:= "SELECT COUNT(*) FROM USERS"

	if err:= config.DB.QueryRow(query).Scan(&count); err != nil {
		log.Println("❌ ERROR: count query failed")
		return 0, err
	}

	return count, nil
}

func DeleteAllUsers(ctx context.Context) error {
	query:= "TRUNCATE TABLE users RESTART IDENTITY"
	
	_, err:= config.DB.ExecContext(ctx, query)
	if err != nil {
		log.Println("Error in deleting User data")
		return err
	}

	log.Println("😉 All User data removed from DB")
	return nil
}

func DeleteUserById(ctx context.Context, userId int) error {
	query:= "DELETE FROM users WHERE id = $1"

	res, err:= config.DB.ExecContext(ctx, query, userId)
	if err != nil {
		log.Println("Error in deleting the user")
		return err
	}

	check, err:= res.RowsAffected()
	if err != nil {
		log.Println("❌ ERROR: In checking rows affected after UPDATE query runs.")
		return err
	}

	if check == 0 {
		log.Println("🚨 0 rows affected, enter a valid ID")
		return errors.New("User not found")
	}

	println("😉 User deleted from DB")
	return nil
}

func UpdateUser(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query:= "UPDATE users SET "
	var args []interface{}
	var parts []string

	i:= 1
	for column, value := range updates {
		parts = append(parts, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	query += strings.Join(parts, ", ") + fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	res, err:= config.DB.Exec(query, args...)
	if err != nil {
		log.Println("❌ ERROR: In updating users")
		return err
	}

	// check does any rows are affected or not (what if user pass an id which is not present in DB)
	check, err:= res.RowsAffected()
	if err != nil {
		log.Println("❌ ERROR: In checking rows affected after UPDATE query runs.")
		return err
	}
	if check == 0 {
		log.Println("🚨 0 rows affected, enter a valid ID")
		return errors.New("User not found")
	}

	log.Println("🥳 User updated successfully in DB")
	return nil
}