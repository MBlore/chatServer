package models

// UserModel is a model of a user in the data store.
type UserModel struct {
	ID          int
	Username    string
	Password    string
	DisplayName *string
	Status      int
	StatusText  *string
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

type PendingContactModel struct {
	ID          int
	Username    string
	DisplayName *string
	ImageURL    *string
	Message     *string
}
