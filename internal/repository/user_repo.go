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
	query:= `
			INSERT INTO users(name, age, password, role_id)
			VALUES($1, $2, $3, $4)
			RETURNING id, role_id
		`
	
	return config.DB.QueryRow(query, user.Name, user.Age, user.Password, 2).Scan(&user.ID, &user.RoleID) // 2 is for by default user role
}

func GetUsers(limit, offset int) ([]models.User, error) {
	rows, err:= config.DB.Query("SELECT id, name, age, role_id FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Println("Get all users query failed to execute")
		return nil, err
	}
	defer rows.Close()

	var userList []models.User

	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Age, &u.RoleID)
		userList = append(userList, u)
	}

	return userList, nil
}

func GetUserById(userId int) (*models.User, error) {
	var u models.User

	err:= config.DB.QueryRow("SELECT id, name, age, role_id FROM users WHERE id = $1", userId).Scan(&u.ID, &u.Name, &u.Age, &u.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with ID %d not found", userId)
			return nil, sql.ErrNoRows
		}
		log.Println("Database error:", err)
		return nil, err
	}

	return &u, nil
}

func GetUserByName(name string) (*models.User, error) {
	// query:= "SELECT id, name, age, password, role FROM users WHERE name=$1"
	query:= `
			SELECT u.id, u.name, u.age, u.password, u.role_id, r.name
			FROM users u
			JOIN roles r ON u.role_id = r.id
			WHERE u.name = $1
		`

	var u models.User
	err:= config.DB.QueryRow(query, name).Scan(&u.ID, &u.Name, &u.Age, &u.Password, &u.RoleID, &u.Role)
	if err != nil {
		log.Printf("User with Name: %s is not fond", name)
		return nil, err
	}

	log.Println("User found with Name: ", name)
	return &u, nil
}

func GetPermissionsByRole(roleId int) ([]string, error) {
	query:= `
			SELECT p.name
			FROM permissions p
			JOIN role_permissions rp
			ON p.id = rp.permission_id
			WHERE rp.role_id = $1
		`

	rows, err:= config.DB.Query(query, roleId)
	if err != nil {
		log.Printf("Permission with this role id %d is not found", roleId)
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			log.Println("❌ ERROR: during copies the columns in the current row")
			return nil, err
		}
		perms = append(perms, perm)
	}

	log.Println("Permission fetched")
	return perms, nil
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

func GetUserAuthInfo(userId int) (*models.User, error) {
	var u models.User

	if err:= config.DB.QueryRow("SELECT role_id FROM users WHERE id=$1", userId).Scan(&u.RoleID); err != nil {
		log.Printf("No role found with this userId: %d", userId)
		return nil, err
	}

	return &u, nil
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
		log.Println("❌ ERROR: In checking rows affected after DELETE query runs.")
		return err
	}

	if check == 0 {
		log.Println("🚨 0 rows affected, enter a valid ID")
		return errors.New("User not found")
	}

	println("😉 User deleted from DB")
	return nil
}

func DeleteUserTx(ctx context.Context, tx *sql.Tx, userId int) error {
	query:= "DELETE FROM users WHERE id = $1"

	res, err:= tx.ExecContext(ctx, query, userId)
	if err != nil {
		log.Println("Error in delete user transaction")
		return err
	}

	check, err:= res.RowsAffected()
	if err != nil {
		log.Println("❌ ERROR: In checking rows affected after DELETE query(transaction) runs.")
		return err
	}

	if check == 0 {
		log.Println("🚨 0 rows affected, enter a valid ID")
		return errors.New("User not found")
	}

	println("😉 User deleted from DB")
	return nil
}

func CreateAuditLogTx(ctx context.Context, tx *sql.Tx, actorId int, action string, targetId int) error {
	query:= `
				INSERT INTO audits(actor_id, action, target_id)
				VALUES($1, $2, $3)
			`

	_, err:= tx.ExecContext(ctx, query, actorId, action, targetId)
	if err != nil {
		log.Println("Insertion failed on Audits")
		return err
	}

	return nil
}

func UpdateUser(id int, updates map[string]interface{}) error {
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