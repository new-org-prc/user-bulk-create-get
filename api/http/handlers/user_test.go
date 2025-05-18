package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"sika/internal/address"
	"sika/internal/user"
	"sika/pkg/storage/entities"
	"sika/service"
	"sika/test/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         string
		mockUser       *entities.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful get user",
			userID: "123",
			mockUser: &entities.User{
				ID:          "123",
				Name:        "Test User",
				Email:       "test@example.com",
				PhoneNumber: "1234567890",
			},
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"id":           "123",
				"name":         "Test User",
				"email":        "test@example.com",
				"phone_number": "1234567890",
				"addresses":    nil,
			},
		},
		{
			name:           "user not found",
			userID:         "999",
			mockUser:       nil,
			mockError:      assert.AnError,
			expectedStatus: fiber.StatusNotFound,
			expectedBody: map[string]interface{}{
				"message": "user not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Fiber app
			app := fiber.New()

			// Create mock repositories
			mockUserRepo := mocks.NewMockUserRepo(ctrl)
			mockAddressRepo := mocks.NewMockAddressRepo(ctrl)

			// Setup mock expectations only if we expect a call
			if tt.userID != "" {
				mockUserRepo.EXPECT().
					GetUserByID(gomock.Any(), tt.userID).
					Return(tt.mockUser, tt.mockError).
					Times(1)
			}

			// Create service with mocks
			userOps := user.NewOps(mockUserRepo)
			addressOps := address.NewOps(mockAddressRepo)
			userService := service.NewUserService(userOps, addressOps)

			// Setup route
			app.Get("/users/:UserID", GetUserByID(userService))

			// Create request
			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response body
			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			require.NoError(t, err)

			// Assert response body
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
