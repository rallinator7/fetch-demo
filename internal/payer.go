package internal

import "github.com/google/uuid"

type Payer struct {
	Name string
	Id   string
}

type PayerAdded struct {
	Payer
}

// NewPayer accepts a name and assigns an Id to the payer. Returns the new payer.
func NewPayer(name string) Payer {
	return Payer{
		Name: name,
		Id:   uuid.NewString(),
	}
}

// NewPayerWithId accepts a name and an id and returns a the new payer.
func NewPayerWithId(id string, name string) Payer {
	return Payer{
		Name: name,
		Id:   uuid.NewString(),
	}
}

// NewPayerAdded accepts a payer and creates a payerAdded event
func NewPayerAdded(payer Payer) PayerAdded {
	return PayerAdded{
		Payer: payer,
	}
}
