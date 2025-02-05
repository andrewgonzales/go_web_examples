package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const insertUser = "INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)"
const queryUser = "SELECT id, username, password, created_at FROM users WHERE id = ?"
const queryUsers = "SELECT id, username, password, created_at FROM users LIMIT 10"
const deleteUser = "DELETE FROM users WHERE id = ?"

func main() {
	// test db in docker container on localhost:3306
	db, err := sql.Open("mysql", "root:abc123@(127.0.0.1:3306)/users?parseTime=true")

	if err != nil {
		fmt.Println("Error opening db")
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging db")
		log.Fatal(err)
	}

	fmt.Printf("db %v \n", db)
	fmt.Printf("err %v \n", err)

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT,
 			username TEXT NOT NULL,
       		password TEXT NOT NULL,
        	created_at DATETIME,
        	PRIMARY KEY (id)
	);
	`

	// Executes the SQL query in our database. Check err to ensure there was no error.
	if _, err := db.Exec(query); err != nil {
		fmt.Println("Error creating user table")
		log.Fatal(err)
	}

	// Insert a user
	// username := "somedude"
	// password := "secret"
	username := "somegal"
	password := "secret2"
	createdAt := time.Now()

	result, err := db.Exec(insertUser, username, password, createdAt)

	if err != nil {
		fmt.Println("Error inserting user")
		log.Fatal(err)
	}

	userId, err := result.LastInsertId()

	if err != nil {
		fmt.Println("Error getting last insert")
		log.Print(err)
	}

	fmt.Println("inserted userId", userId)

	{ // Query a single user
		var (
			id        int
			username  string
			password  string
			createdAt time.Time
		)

		if err := db.QueryRow(queryUser, 1).Scan(&id, &username, &password, &createdAt); err != nil {
			fmt.Println("Error querying user")
			log.Print(err)
		}

		fmt.Println("User query: ", id, username, password, createdAt)
	}

	{ // Query all users
		type user struct {
			id        int
			username  string
			password  string
			createdAt time.Time
		}

		rows, err := db.Query(queryUsers)
		if err != nil {
			fmt.Println("Error querying users")
			log.Fatal(err)
		}
		defer rows.Close()

		var users []user
		for rows.Next() {
			var u user

			err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
			if err != nil {
				fmt.Println("Error scanning user")
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			fmt.Println("Error iterating users")
			log.Fatal(err)
		}

		fmt.Printf("users %#v", users)
	}

	{ // Delete a user
		_, err := db.Exec(deleteUser, 1)
		if err != nil {
			fmt.Println("Error deleting user")
			log.Print(err)
		}
	}

}
