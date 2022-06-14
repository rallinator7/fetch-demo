package user

import "fmt"

type MissingUserError struct {
	User string
}

func (e *MissingUserError) Error() string {
	return fmt.Sprintf("user %s does not exist", e.User)
}

type DuplicateUserError struct {
	User string
}

func (e *DuplicateUserError) Error() string {
	return fmt.Sprintf("user %s already exists", e.User)
}

type NegativeTransactionError struct {
	Target string
}

func (e *NegativeTransactionError) Error() string {
	return fmt.Sprintf("total points for %s cannot be negative", e.Target)
}
