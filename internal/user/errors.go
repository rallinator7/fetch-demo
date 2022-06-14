package user

import (
	"fmt"
)

type DuplicateUserError struct {
	First string
	Last  string
}

func (e *DuplicateUserError) Error() string {
	return fmt.Sprintf("a user with the name %s %s already exists", e.First, e.Last)
}
