package service

import (
	"context"
	"testing"
	"time"

	"sika/internal/address"
	"sika/internal/user"
	"sika/pkg/load"
	"sika/pkg/storage/entities"
	"sika/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_ImportUsers(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockAddressRepo := mocks.NewMockAddressRepo(ctrl)

	userOps := user.NewOps(mockUserRepo)
	addressOps := address.NewOps(mockAddressRepo)

	service := NewUserService(userOps, addressOps)

	tests := []struct {
		name       string
		usersData  []load.User
		setupMocks func()
		wantErr    bool
	}{
		{
			name: "successful import",
			usersData: []load.User{
				{
					ID:          "1",
					Name:        "Test User",
					Email:       "test@example.com",
					PhoneNumber: "1234567890",
					Addresses: []load.Address{
						{
							Street:  "123 Test St",
							City:    "Test City",
							State:   "Test State",
							ZipCode: "12345",
							Country: "Test Country",
						},
					},
				},
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
				mockAddressRepo.EXPECT().CreateBatchAddresses(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user creation fails",
			usersData: []load.User{
				{
					ID:          "1",
					Name:        "Test User",
					Email:       "test@example.com",
					PhoneNumber: "1234567890",
				},
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "address creation fails",
			usersData: []load.User{
				{
					ID:          "1",
					Name:        "Test User",
					Email:       "test@example.com",
					PhoneNumber: "1234567890",
					Addresses: []load.Address{
						{
							Street:  "123 Test St",
							City:    "Test City",
							State:   "Test State",
							ZipCode: "12345",
							Country: "Test Country",
						},
					},
				},
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
				mockAddressRepo.EXPECT().CreateBatchAddresses(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Execute
			err := service.ImportUsers(tt.usersData)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockAddressRepo := mocks.NewMockAddressRepo(ctrl)

	userOps := user.NewOps(mockUserRepo)
	addressOps := address.NewOps(mockAddressRepo)

	service := NewUserService(userOps, addressOps)

	tests := []struct {
		name       string
		userID     string
		setupMocks func()
		want       *entities.User
		wantErr    bool
	}{
		{
			name:   "successful get",
			userID: "1",
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByID(gomock.Any(), "1").Return(&entities.User{
					ID:          "1",
					Name:        "Test User",
					Email:       "test@example.com",
					PhoneNumber: "1234567890",
				}, nil)
			},
			want: &entities.User{
				ID:          "1",
				Name:        "Test User",
				Email:       "test@example.com",
				PhoneNumber: "1234567890",
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userID: "2",
			setupMocks: func() {
				mockUserRepo.EXPECT().GetUserByID(gomock.Any(), "2").Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Execute
			got, err := service.GetUserByID(context.Background(), tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWorkerPool_ConcurrentProcessing(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockAddressRepo := mocks.NewMockAddressRepo(ctrl)

	userOps := user.NewOps(mockUserRepo)
	addressOps := address.NewOps(mockAddressRepo)

	service := NewUserService(userOps, addressOps)

	// Create test data
	usersData := make([]load.User, 100)
	for i := 0; i < 100; i++ {
		usersData[i] = load.User{
			ID:          string(rune(i + 1)),
			Name:        "Test User",
			Email:       "test@example.com",
			PhoneNumber: "1234567890",
			Addresses: []load.Address{
				{
					Street:  "123 Test St",
					City:    "Test City",
					State:   "Test State",
					ZipCode: "12345",
					Country: "Test Country",
				},
			},
		}
	}

	// Setup mocks
	mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(100)
	mockAddressRepo.EXPECT().CreateBatchAddresses(gomock.Any(), gomock.Any()).Return(nil).Times(100)

	// Execute with timeout
	done := make(chan error)
	go func() {
		done <- service.ImportUsers(usersData)
	}()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Import operation timed out")
	}
}
