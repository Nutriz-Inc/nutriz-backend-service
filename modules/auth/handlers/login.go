package handlers

import (
	c "context"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/auth/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type HandlerLogin struct {
	userRepo    *repositories.UserRepository
	addressRepo *repositories.AddressRepository
	config      *config.Env
}

func HandlerLoginStart(userRepo *repositories.UserRepository, addressRepo *repositories.AddressRepository, config *config.Env) *HandlerLogin {
	return &HandlerLogin{userRepo, addressRepo, config}
}

func (h *HandlerLogin) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.LoginReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerLogin) Execute(ctx c.Context, data *dto.LoginReq) (*dto.LoginRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorBadRequest("Invalid email or password", "auth.invalid_credentials")
	}

	secret := utils.Secret{}

	isPasswordValid := secret.IsEqual(h.config, data.Password, user.Password)
	if !isPasswordValid {
		return nil, fluxgo.ErrorBadRequest("Invalid email or password", "auth.invalid_credentials")
	}

	const SEVEN_DAYS = 7 * 24 * time.Hour

	tokenPayload := utils.JwtClaims{
		IdUser: user.IdUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(SEVEN_DAYS)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := utils.CreateJwtToken(h.config, tokenPayload)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create jwt token")
	}

	isRecurrentDonor := user.IsBloodExamValid()

	return &dto.LoginRes{
		Token:            token,
		IdUser:           user.IdUser,
		Name:             user.Name,
		Type:             user.Type,
		IsRecurrentDonor: &isRecurrentDonor,
	}, nil
}
