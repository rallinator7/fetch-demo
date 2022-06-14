package payer

import "fmt"

type MissingUserError struct {
	Payer string
}

func (e *MissingUserError) Error() string {
	return fmt.Sprintf("Payer %s does not exist", e.Payer)
}

type DuplicateUserError struct {
	Payer string
}

func (e *DuplicateUserError) Error() string {
	return fmt.Sprintf("Payer %s already exists", e.Payer)
}
