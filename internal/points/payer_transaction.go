package points

import "time"

type PayerPoints struct {
	Payer  string
	Amount int
}

type PayerTransaction struct {
	TimeStamp time.Time
	PayerPoints
}

func NewPayerPoints(payer string, amount int) PayerPoints {
	return PayerPoints{
		Payer:  payer,
		Amount: amount,
	}
}

func NewPayerTransaction(payer string, amount int, timestamp string) (PayerTransaction, error) {
	stamp, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return PayerTransaction{}, err
	}

	transaction := PayerTransaction{
		PayerPoints: PayerPoints{
			Payer:  payer,
			Amount: amount,
		},
		TimeStamp: stamp,
	}

	return transaction, nil
}
