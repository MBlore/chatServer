package dbaccess

import (
	"chatServer/models"
	"database/sql"
	"errors"
	"log"

	// We use the mysql driver.
	_ "github.com/go-sql-driver/mysql"
)

// Status ID of a user.
const (
	StatusOffline = 0
	StatusOnline  = 1
)

var database *sql.DB

// OpenConnection opens a connection to the database.
func OpenConnection() {
	db, err := sql.Open("mysql", "")
	if err != nil {
		panic(err.Error())
	}

	database = db
}

// CloseConnection closes the database connection.
func CloseConnection() {
	database.Close()
}

// GetUserByUsername returns a user model from the data store by username.
func GetUserByUsername(username string) (*models.UserModel, error) {
	rows, err := database.Query("call getUserByUsername(?)", username)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		user := models.UserModel{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.DisplayName, &user.Status, &user.StatusText)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, nil
}

// GetUserByID returns a user model from the data store by id.
func GetUserByID(userID int) (*models.UserModel, error) {
	rows, err := database.Query("call getUserById(?)", userID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		user := models.UserModel{}
		err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Status, &user.ImageURL, &user.StatusText)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, nil
}

// LoginUser logs a user in and sets their status to online.
func LoginUser(userID int) error {
	_, err := database.Exec("call loginUser(?)", userID)
	if err != nil {
		return err
	}

	return nil
}

// LogoutUser logs a user out and sets their status to offline.
func LogoutUser(userID int) error {
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
func GetFriends(userID int) ([]models.FriendModel, error) {
	rows, err := database.Query("call getFriends(?)", userID)
	if err != nil {
		return nil, err
	}

	friends := []models.FriendModel{}

	for rows.Next() {
		friend := models.FriendModel{}

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
func ResetUserStatuses() error {
	_, err := database.Exec("call resetUserStatuses()")
	if err != nil {
		return err
	}

	return nil
}

// CreateAccount creates a new account.
func CreateAccount(username string, password string, email string, displayName string, validationGUID string) error {
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

// AddPendingContact adds a contact request for the requested user.
func AddPendingContact(userID int, userAddingID int, message *string) error {
	res, err := database.Exec(
		"call addPendingContact(?,?,?)",
		userID,
		userAddingID,
		message)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return errors.New("failed to add pending contact")
	}

	return nil
}

// GetPendingContact retreives a pending contact by user and requested user.
func GetPendingContact(requestedUserID int, addingUserID int) (*int, error) {
	row := database.QueryRow("call getPendingContact(?,?)", requestedUserID, addingUserID)

	var rowID int

	err := row.Scan(&rowID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &rowID, nil
}

// GetUserContactByContactUserID retreives a users contact.
func GetUserContactByContactUserID(userID int, contactUserID int) (*int, error) {
	row := database.QueryRow("call getUserContactByContactUserID(?,?)", userID, contactUserID)

	var rowID int

	err := row.Scan(&rowID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &rowID, nil
}

// GetUserPendingContacts returns a list of a users pending contact requests.
func GetUserPendingContacts(userID int) ([]models.PendingContactModel, error) {
	rows, err := database.Query("call getUserPendingContacts(?)", userID)
	if err != nil {
		return nil, err
	}

	contacts := []models.PendingContactModel{}

	for rows.Next() {
		c := models.PendingContactModel{}

		err := rows.Scan(&c.ID, &c.Username, &c.DisplayName, &c.ImageURL, &c.Message)
		if err != nil {
			log.Printf("Failed to map pending contact model.")
			continue
		}

		contacts = append(contacts, c)
	}

	return contacts, nil
}

// ConfirmContactRequest confirms a contact request.
func ConfirmContactRequest(requestedUserID int, addingUserID int) error {
	res, err := database.Exec(
		"call confirmContact(?,?)",
		requestedUserID,
		addingUserID)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return errors.New("failed to confirm pending contact")
	}

	return nil
}

// RejectContactRequest rejects a contact request.
func RejectContactRequest(requestedUserID int, addingUserID int) error {
	res, err := database.Exec(
		"call rejectContact(?,?)",
		requestedUserID,
		addingUserID)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return errors.New("failed to reject pending contact")
	}

	return nil
}

// SetStatus sets a users status.
func SetStatus(userID int, statusID int) error {
	res, err := database.Exec(
		"call setStatus(?,?)",
		userID,
		statusID)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra != 1 {
		return errors.New("failed to set user status")
	}

	return nil
}
