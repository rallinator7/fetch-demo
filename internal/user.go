package internal

import "github.com/google/uuid"

type User struct {
	Id        string
	FirstName string
	LastName  string
}

type UserAdded struct {
	User
}

// NewUser accepts a first name and last name and assigns an id.  Returns the new user.
func NewUser(first string, last string) User {
	return User{
		FirstName: first,
		LastName:  last,
		Id:        uuid.NewString(),
	}
}

// NewUserWithID accepts all params as NewUser as well as an id.  Assigns params to a user and returns it.
func NewUserWithId(id string, first string, last string) User {
	return User{
		FirstName: first,
		LastName:  last,
		Id:        id,
	}
}

// NewUserAdded accepts all params as NewUser as well as an id.  Assigns params to a user and returns it.
func NewUserAdded(user User) UserAdded {
	return UserAdded{
		User: user,
	}
}
