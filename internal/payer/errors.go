package payer

import (
	"fmt"
)

type DuplicatePayerError struct {
	Name string
}

func (e *DuplicatePayerError) Error() string {
	return fmt.Sprintf("a payer with the name %s already exists", e.Name)
}
