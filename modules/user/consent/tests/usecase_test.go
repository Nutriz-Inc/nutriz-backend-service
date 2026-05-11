package tests

import (
	"context"
	"errors"
	"nutriz-backend-service/modules/user/consent/dtos"
	"nutriz-backend-service/modules/user/consent/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConsentRepo struct {
	mock.Mock
}

func (m *mockConsentRepo) CreateConsent(
	ctx context.Context,
	data dtos.CreateConsentReq,
	idUser string,
	idConsentLog string,
) error {
	args := m.Called(ctx, data, idUser, idConsentLog)
	return args.Error(0)
}

func TestCreateConsentUseCase_Execute_Success(t *testing.T) {
	repo := new(mockConsentRepo)
	uc := usecases.NewCreateConsentUseCase(repo)

	req := &dtos.CreateConsentReq{
		TermsVersion: "v1.0",
		IpAddress:    "192.168.0.1",
		UserAgent:    "Mozilla/5.0",
	}

	repo.On(
		"CreateConsent",
		mock.Anything,                 // ctx
		*req,                          // data
		"user-123",                    // idUser
		mock.AnythingOfType("string"), // idConsentLog (uuid gerado internamente)
	).Return(nil)

	resp, globErr := uc.Execute(context.Background(), "user-123", req)

	assert.Nil(t, globErr)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-123", resp.IdUser)
	assert.Equal(t, "v1.0", resp.TermsVersion)
	assert.NotEmpty(t, resp.IdConsentLog, "IDConsentLog deve ser gerado pelo servidor")

	repo.AssertExpectations(t)
}

func TestCreateConsentUseCase_Execute_RepositoryError(t *testing.T) {
	repo := new(mockConsentRepo)
	uc := usecases.NewCreateConsentUseCase(repo)

	req := &dtos.CreateConsentReq{
		TermsVersion: "v1.0",
		IpAddress:    "192.168.0.1",
		UserAgent:    "Mozilla/5.0",
	}

	repo.On(
		"CreateConsent",
		mock.Anything,
		*req,
		"user-123",
		mock.AnythingOfType("string"),
	).Return(errors.New("db: connection refused"))

	resp, globErr := uc.Execute(context.Background(), "user-123", req)

	assert.Nil(t, resp)
	assert.NotNil(t, globErr)

	repo.AssertExpectations(t)
}
