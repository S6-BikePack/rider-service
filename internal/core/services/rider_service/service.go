package rider_service

import (
	"errors"
	"github.com/google/uuid"
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

func (srv *service) Get(uuid uuid.UUID) (domain.Rider, error) {
	return srv.riderRepository.Get(uuid)
}

func (srv *service) Create(name string, status int8) (domain.Rider, error) {
	rider := domain.NewRider(name, status, domain.Location{})

	rider, err := srv.riderRepository.Save(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.CreateRider(rider)
	return rider, nil
}

func (srv *service) Update(uuid uuid.UUID, name string, status int8) (domain.Rider, error) {
	rider, err := srv.Get(uuid)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.Status = status

	if name != "" {
		rider.Name = name
	}

	rider, err = srv.riderRepository.Update(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.UpdateRider(rider)
	return rider, nil
}

func (srv *service) UpdateLocation(uuid uuid.UUID, location domain.Location) (domain.Rider, error) {
	rider, err := srv.Get(uuid)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.Location = location

	rider, err = srv.riderRepository.Update(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	srv.messagePublisher.UpdateRider(rider)
	return rider, nil
}
