package rider_service

import (
	"errors"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/ports"
)

type service struct {
	riderRepository  ports.RiderRepository
	messagePublisher ports.MessageBusPublisher
}

func New(riderRepository ports.RiderRepository, messagePublisher ports.MessageBusPublisher) *service {
	return &service{
		riderRepository:  riderRepository,
		messagePublisher: messagePublisher,
	}
}

func (srv *service) GetAll() ([]domain.Rider, error) {
	return srv.riderRepository.GetAll()
}

func (srv *service) Get(id string) (domain.Rider, error) {
	return srv.riderRepository.Get(id)
}

func (srv *service) Create(userId string, status int8) (domain.Rider, error) {
	user, err := srv.riderRepository.GetUser(userId)

	rider := domain.NewRider(user, status, domain.Location{})

	rider, err = srv.riderRepository.Save(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.CreateRider(rider)
	return rider, nil
}

func (srv *service) Update(id string, status int8) (domain.Rider, error) {
	rider, err := srv.Get(id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.Status = status

	rider, err = srv.riderRepository.Update(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.UpdateRider(rider)
	return rider, nil
}

func (srv *service) UpdateLocation(id string, location domain.Location) (domain.Rider, error) {
	rider, err := srv.Get(id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.Location = location

	rider, err = srv.riderRepository.Update(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.UpdateRiderLocation(rider.UserID, location)
	return rider, nil
}

func (srv *service) SaveOrUpdateUser(user domain.User) error {
	if user.Name == "" || user.LastName == "" || user.ID == "" {
		return errors.New("incomplete user data")
	}

	err := srv.riderRepository.SaveOrUpdateUser(user)

	return err
}
