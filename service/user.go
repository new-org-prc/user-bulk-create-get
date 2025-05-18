package service

import (
	"context"
	"fmt"
	"sika/internal/address"
	"sika/internal/user"
	"sika/pkg/load"
	"sika/pkg/storage/entities"
	"sync"
)

type UserService struct {
	userOps    *user.Ops
	addressOps *address.Ops
}

func NewUserService(userOps *user.Ops, addressOps *address.Ops) *UserService {
	return &UserService{
		userOps:    userOps,
		addressOps: addressOps,
	}
}

type Job struct {
	user      *entities.User
	addresses []*entities.Address
}

type WorkerPool struct {
	numOfWorkers int
	wg           sync.WaitGroup
	jobs         chan Job
	results      chan error
}

func NewWorkerPool(n int) *WorkerPool {
	return &WorkerPool{
		numOfWorkers: n,
		wg:           sync.WaitGroup{},
		jobs:         make(chan Job, 1000000),
		results:      make(chan error, 1000000),
	}
}

func (s *UserService) ImportUsers(usersData []load.User) error {
	wp := NewWorkerPool(10)
	ctx := context.Background()
	for i := 0; i < wp.numOfWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for j := range wp.jobs {
				err := s.userOps.CreateUser(ctx, j.user)
				if err != nil {
					wp.results <- fmt.Errorf("user with Id %s insertion to db failed %w", j.user.ID, err)
					continue
				}

				if len(j.addresses) > 0 {
					addresses := make([]entities.Address, len(j.addresses))
					for i, a := range j.addresses {
						addresses[i] = *a
					}

					err := s.addressOps.CreateBatchAddress(ctx, addresses)
					if err != nil {
						wp.results <- fmt.Errorf("batch address insertion for userID %s failed %w", j.user.ID, err)
						continue
					}
				}
				wp.results <- nil
			}
		}()
	}

	go func() {
		for _, u := range usersData {
			uEntity := &entities.User{
				ID:          u.ID,
				Name:        u.Name,
				Email:       u.Email,
				PhoneNumber: u.PhoneNumber,
			}
			var eAddresses []*entities.Address
			for _, a := range u.Addresses {
				aEntity := &entities.Address{
					UserID:  u.ID,
					Street:  a.Street,
					City:    a.City,
					State:   a.State,
					ZipCode: a.ZipCode,
					Country: a.Country,
				}
				eAddresses = append(eAddresses, aEntity)
			}

			j := Job{
				user:      uEntity,
				addresses: eAddresses,
			}

			wp.jobs <- j

		}
		close(wp.jobs)
	}()
	var errs []error
	for i := 0; i < len(usersData); i++ {
		if v, ok := <-wp.results; ok {
			if v != nil {
				errs = append(errs, v)
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during import %v", len(errs), errs)
	}
	wp.wg.Wait()
	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.userOps.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %s: %w", id, err)
	}
	return user, nil
}

func (s *UserService) ClearUserAndAddressDataFromDB() error {
	err := s.addressOps.ClearAllAddressesDataFromDB()
	if err != nil {
		return fmt.Errorf("failed to clear addresses table: %w", err)
	}
	err = s.userOps.ClearAllUsersDataFromDB()
	if err != nil {
		return fmt.Errorf("failed to clear users table: %w", err)
	}
	return nil
}
