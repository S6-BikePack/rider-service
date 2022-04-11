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
	rider, err := srv.riderRepository.Get(id)

	if err != nil {
		return rider, err
	}

	if (rider == domain.Rider{}) {
		return rider, errors.New("could not find rider with id: " + id)
	}

	return rider, nil
}

func (srv *service) Create(userId string, serviceArea int, capacity domain.Dimensions) (domain.Rider, error) {
	user, err := srv.riderRepository.GetUser(userId)

	rider := domain.NewRider(user, 0, serviceArea, capacity)

	rider, err = srv.riderRepository.Save(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.CreateRider(rider)
	return rider, nil
}

func (srv *service) Update(id string, status int, serviceArea int, capacity domain.Dimensions) (domain.Rider, error) {
	rider, err := srv.Get(id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.ServiceArea = serviceArea
	rider.Status = status
	rider.Capacity = capacity

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
