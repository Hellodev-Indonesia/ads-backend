package auth_test

import (
	"testing"

	"github.com/alex/ads_backend/internal/core/auth"
	"github.com/alex/ads_backend/internal/core/auth/dto"
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupService(t *testing.T) (*user.MockRepository, *permission.MockRepository, auth.Service) {
	mockUserRepo := user.NewMockRepository(t)
	mockPermRepo := permission.NewMockRepository(t)
	svc := auth.NewService(mockUserRepo, mockPermRepo)
	return mockUserRepo, mockPermRepo, svc
}

func TestService_Login(t *testing.T) {
	mockUserRepo, _, svc := setupService(t)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &user.User{
		ID:       1,
		Email:    "test@test.com",
		Password: string(hashedPassword),
	}

	mockUserRepo.On("FindByEmail", "test@test.com").Return(existingUser, nil)

	req := dto.LoginRequest{
		Email:    "test@test.com",
		Password: "password123",
	}

	resp, err := svc.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, uint(1), resp.User.ID)
	assert.Equal(t, "test@test.com", resp.User.Email)
}

func TestService_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo, _, svc := setupService(t)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &user.User{
		ID:       1,
		Email:    "test@test.com",
		Password: string(hashedPassword),
	}

	mockUserRepo.On("FindByEmail", "test@test.com").Return(existingUser, nil)

	req := dto.LoginRequest{
		Email:    "test@test.com",
		Password: "wrongpassword",
	}

	resp, err := svc.Login(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid email or password", err.Error())
}
