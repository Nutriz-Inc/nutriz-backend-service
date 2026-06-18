package handlers

import (
	c "context"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerUpdateUser struct {
	config   *config.Env
	userRepo *repositories.UserRepository
}

func HandlerUpdateUserStart(config *config.Env, userRepo *repositories.UserRepository) *HandlerUpdateUser {
	return &HandlerUpdateUser{config, userRepo}
}

func (h *HandlerUpdateUser) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateUserReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateUser) Execute(ctx c.Context, data *dto.UpdateUserReq) (*dto.UpdateUserRes, *fluxgo.GlobalError) {
	//busca o user alvo
	user, err := h.userRepo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	// busca que m tá fazendo a ação
	actionBy, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get action by user")
	}
	if actionBy == nil {
		return nil, fluxgo.ErrorNotFound("Action by user not found")
	}

	// só o proprio user editar
	if actionBy.IdUser != user.IdUser {
		return nil, utils.ErrorForbidden("You don't have permission to update this user", "user.forbidden")
	}

	// type so pode ser trocado por adm
	if data.Type != nil && actionBy.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden("Only admins can change user type", "user.forbidden_type")
	}

	// common não recebe internal_identifier
	if data.InternalIdentifier != nil && user.Type == entities.EnumUserTypeCommon {
		return nil, fluxgo.ErrorBadRequest("Common users cannot have internal identifier", "user.invalid_field")
	}

	// valida email único
	if data.Email != nil {
		userWithSameEmail, err := h.userRepo.GetUserByEmail(ctx, *data.Email)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get user by email")
		}
		if userWithSameEmail != nil && userWithSameEmail.IdUser != user.IdUser {
			return nil, fluxgo.ErrorBadRequest("Email already in use", "user.duplicate_email")
		}
	}

	// hash da nova senha
	var hashPassword *string
	if data.Password != nil {
		secret := utils.Secret{}
		hash, err := secret.Encrypt(h.config, *data.Password)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to encrypt password")
		}
		hashPassword = &hash
	}

	err = h.userRepo.UpdateUser(ctx, &repositories.UpdateUserRepositoryReq{
		IdUser:             data.Id,
		ActionBy:           data.ActionBy,
		InternalIdentifier: data.InternalIdentifier,
		Type:               data.Type,
		Name:               data.Name,
		PhoneNumber:        data.PhoneNumber,
		Email:              data.Email,
		Password:           hashPassword,
	})
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update user")
	}

	updatedUser, err := h.userRepo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get updated user")
	}

	return &dto.UpdateUserRes{User: *updatedUser}, nil
}	