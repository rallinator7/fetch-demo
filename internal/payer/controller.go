package payer

import (
	internal "github.com/rallinator7/fetch-demo/internal"
)

type PayerRepoer interface {
	Add(payer internal.Payer) error
}

type PayerPublisher interface {
	PlayerAdded(payer internal.Payer) error
}

type Logger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type Controller struct {
	repo      PayerRepoer
	publisher PayerPublisher
	logger    Logger
}

func NewController(repo PayerRepoer, publisher PayerPublisher, logger Logger) Controller {
	return Controller{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// AddPayer calls the repository to save a payer to memory.
func (controller *Controller) AddPayer(payer internal.Payer) (internal.Payer, error) {
	// needs work to be atomic but good enough for demo
	err := controller.repo.Add(payer)
	if err != nil {
		controller.logger.Errorw("failed adding payer to repo", "error", err)
		return internal.Payer{}, err
	}

	err = controller.publisher.PlayerAdded(payer)
	if err != nil {
		controller.logger.Errorw("failed publishing payer", "error", err)
		return internal.Payer{}, err
	}

	controller.logger.Infow(
		"payer added",
		"id", payer.Id,
		"name", payer.Name,
	)

	return payer, nil
}
