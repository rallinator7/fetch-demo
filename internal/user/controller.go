package user

import (
	"fmt"

	internal "github.com/rallinator7/fetch-demo/internal"
)

type Logger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type UserPublisher interface {
	UserAdded(user internal.User) error
}

type UserRepoer interface {
	Add(user internal.User) error
}

type Controller struct {
	repo      UserRepoer
	publisher UserPublisher
	logger    Logger
}

func NewController(repo UserRepoer, publisher UserPublisher, logger Logger) Controller {
	return Controller{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// AddUser calls the repository to save a user to memory.
func (controller *Controller) AddUser(user internal.User) (internal.User, error) {
	// needs work to be atomic but good enough for demo
	err := controller.repo.Add(user)
	if err != nil {
		controller.logger.Errorw("failed adding user to repo", "error", err)
		return internal.User{}, err
	}

	err = controller.publisher.UserAdded(user)
	if err != nil {
		controller.logger.Errorw("failed publishing user", "error", err)
		return internal.User{}, err
	}

	controller.logger.Infow(
		"user added",
		"id", user.Id,
		"name", fmt.Sprintf("%s %s", user.FirstName, user.LastName),
	)

	return user, nil
}
