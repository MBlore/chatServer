package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Status ID of a user.
const (
	StatusOffline = 0
	StatusOnline  = 1
)

// UserModel is a model of a user in the data store.
type UserModel struct {
	ID          int
	Username    string
	Password    string
	DisplayName *string
	Status      int
}

// FriendModel is a model of users friend.
type FriendModel struct {
	ID          int
	Username    string
	DisplayName *string
	Status      int
	ImageURL    *string
	StatusText  *string
}

type dbAccess struct{}

// DBAccess provides database access functions.
var DBAccess dbAccess

var database *sql.DB

// OpenConnection opens a connection to the database.
func (dbAccess) OpenConnection() {
	db, err := sql.Open("mysql", "")
	if err != nil {
		panic(err.Error())
	}

	database = db
}

// CloseConnection closes the database connection.
func (dbAccess) CloseConnection() {
	database.Close()
}

// GetUserByUsername returns a user model from the data store by username.
func (dbAccess) GetUserByUsername(username string) (*UserModel, error) {
	rows, err := database.Query("call getUserByUsername(?)", username)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		user := UserModel{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.DisplayName, &user.Status)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, nil
}

// LoginUser logs a user in and sets their status to online.
func (dbAccess) LoginUser(userID int) error {
	_, err := database.Exec("call loginUser(?)", userID)
	if err != nil {
		return err
	}

	return nil
}

// LogoutUser logs a user out and sets their status to offline.
func (dbAccess) LogoutUser(userID int) error {
	res, err := database.Exec("call logoffUser(?)", userID)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return errors.New("failed to update")
	}

	return nil
}

// GetFriends returns a list of friends that the specified user has.
func (dbAccess) GetFriends(userID int) ([]FriendModel, error) {
	rows, err := database.Query("call getFriends(?)", userID)
	if err != nil {
		return nil, err
	}

	friends := []FriendModel{}

	for rows.Next() {
		friend := FriendModel{}

		err := rows.Scan(&friend.ID, &friend.Username, &friend.DisplayName, &friend.Status, &friend.ImageURL, &friend.StatusText)
		if err != nil {
			log.Printf("Failed to map friend model.")
			continue
		}

		friends = append(friends, friend)
	}

	return friends, nil
}

// ResetUserStatuses sets all users to be offline.
func (dbAccess) ResetUserStatuses() error {
	_, err := database.Exec("call resetUserStatuses()")
	if err != nil {
		return err
	}

	return nil
}

func (dbAccess) CreateAccount(username string, password string, email string, displayName string, validationGUID string) error {
	res, err := database.Exec(
		"call createAccount(?,?,?,?,?)",
		username,
		email,
		password,
		displayName,
		validationGUID)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return errors.New("failed to create")
	}

	return nil
}
