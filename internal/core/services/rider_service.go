package services

import (
	"context"
	"errors"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
)

type riderService struct {
	riderRepository  interfaces.RiderRepository
	messagePublisher interfaces.MessageBusPublisher
}

func NewRiderService(riderRepository interfaces.RiderRepository, messagePublisher interfaces.MessageBusPublisher) *riderService {
	return &riderService{
		riderRepository:  riderRepository,
		messagePublisher: messagePublisher,
	}
}

func (srv *riderService) GetAll(ctx context.Context) ([]domain.Rider, error) {
	return srv.riderRepository.GetAll(ctx)
}

func (srv *riderService) Get(ctx context.Context, id string) (domain.Rider, error) {
	rider, err := srv.riderRepository.Get(ctx, id)

	if err != nil {
		return rider, err
	}

	if (rider == domain.Rider{}) {
		return rider, errors.New("could not find rider with id: " + id)
	}

	return rider, nil
}

func (srv *riderService) Create(ctx context.Context, userId string, serviceArea int, capacity domain.Dimensions) (domain.Rider, error) {
	user, err := srv.riderRepository.GetUser(ctx, userId)

	if err != nil {
		return domain.Rider{}, err
	}

	rider := domain.NewRider(user, 0, serviceArea, capacity)

	rider, err = srv.riderRepository.Save(ctx, rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	_ = srv.messagePublisher.CreateRider(ctx, rider)

	return rider, nil
}

func (srv *riderService) Update(ctx context.Context, id string, status int, serviceArea int, capacity domain.Dimensions) (domain.Rider, error) {
	rider, err := srv.Get(ctx, id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.ServiceAreaID = serviceArea
	rider.Status = status

	if capacity != (domain.Dimensions{}) {
		rider.Capacity = capacity
	}

	rider, err = srv.riderRepository.Update(ctx, rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	_ = srv.messagePublisher.UpdateRider(ctx, rider)

	return rider, nil
}

func (srv *riderService) UpdateLocation(ctx context.Context, id string, location domain.Location) (domain.Rider, error) {
	rider, err := srv.Get(ctx, id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.Location = location

	rider, err = srv.riderRepository.Update(ctx, rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	err = srv.messagePublisher.UpdateRiderLocation(ctx, rider.ServiceArea, rider.UserID, location)

	if err != nil {
		return rider, err
	}

	return rider, nil
}

func (srv *riderService) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	if user.Name == "" || user.LastName == "" || user.ID == "" {
		return errors.New("incomplete user data")
	}

	err := srv.riderRepository.SaveOrUpdateUser(ctx, user)

	return err
}
