package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUser(db *sql.DB, user User) error {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at"
	err := db.QueryRow(query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func GetUserById(db *sql.DB, id int) (User, error) {
	var user User

	query := "SELECT id, name, email, created_at FROM users WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return user, err
	}

	return user, nil
}

func GetUsers(db *sql.DB) ([]User, error) {

	rows, err := db.Query("SELECT id, name, email, created_at FROM users;")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		users = append(users, user)
		fmt.Println(users)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return nil, err
	}

	return users, nil
}

func UpdateUser(db *sql.DB, id int, user User) (User, error) {
	fmt.Printf("updating user id: %d for user %v", id, user)
	query := `
	UPDATE users 
	SET name = $1, email = $2, updated_at = NOW() 
	WHERE id = $3 
	RETURNING id, name, email, created_at, updated_at`

	var updatedUser User
	err := db.QueryRow(query, user.Name, user.Email, id).
		Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)

	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func DeleteUser(db *sql.DB, id int) error {

	query := "DELETE FROM users WHERE id = $1"
	res, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil

}
