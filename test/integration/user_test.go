package integration

import (
	"context"
	"testing"

	"sika/internal/address"
	"sika/internal/user"
	"sika/pkg/load"
	"sika/pkg/storage"
	"sika/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Integration(t *testing.T) {
	// Setup test database
	db := SetupTestDB(t)
	defer db.Close()

	// Create repositories
	userRepo := storage.NewUserRepo((*db).DB)
	addressRepo := storage.NewAddressRepo((*db).DB)

	// Create operations
	userOps := user.NewOps(userRepo)
	addressOps := address.NewOps(addressRepo)

	// Create service
	userService := service.NewUserService(userOps, addressOps)

	t.Run("Import and Retrieve User", func(t *testing.T) {
		// Clean up before test
		db.Cleanup(t)

		// Test data
		testUser := load.User{
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
		}

		// Import user
		err := userService.ImportUsers([]load.User{testUser})
		require.NoError(t, err)

		// Retrieve user
		user, err := userService.GetUserByID(context.Background(), testUser.ID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Name, user.Name)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.PhoneNumber, user.PhoneNumber)
	})

	t.Run("Import Multiple Users", func(t *testing.T) {
		// Clean up before test
		db.Cleanup(t)

		// Test data
		users := make([]load.User, 10)
		for i := 0; i < 10; i++ {
			users[i] = load.User{
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

		// Import users
		err := userService.ImportUsers(users)
		require.NoError(t, err)

		// Verify all users were imported
		for _, testUser := range users {
			user, err := userService.GetUserByID(context.Background(), testUser.ID)
			require.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, testUser.ID, user.ID)
		}
	})

	t.Run("Get Non-existent User", func(t *testing.T) {
		// Clean up before test
		db.Cleanup(t)

		// Try to get non-existent user
		user, err := userService.GetUserByID(context.Background(), "non-existent")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_ConcurrentImport(t *testing.T) {
	// Setup test database
	db := SetupTestDB(t)
	defer db.Close()

	// Create repositories
	userRepo := storage.NewUserRepo((*db).DB)
	addressRepo := storage.NewAddressRepo((*db).DB)

	// Create operations
	userOps := user.NewOps(userRepo)
	addressOps := address.NewOps(addressRepo)

	// Create service
	userService := service.NewUserService(userOps, addressOps)

	// Clean up before test
	db.Cleanup(t)

	// Create test data
	users := make([]load.User, 100)
	for i := 0; i < 100; i++ {
		users[i] = load.User{
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

	// Import users
	err := userService.ImportUsers(users)
	require.NoError(t, err)

	// Verify all users were imported
	for _, testUser := range users {
		user, err := userService.GetUserByID(context.Background(), testUser.ID)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
	}
}
