package handlers

import (
	c "context"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
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
	return &HandlerUpdateUser{
		config:   config,
		userRepo: userRepo,
	}
}

func (h *HandlerUpdateUser) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateUserReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateUser) Execute(ctx c.Context, data *dto.UpdateUserReq) (*dto.UpdateUserRes, *fluxgo.GlobalError) {
	actor, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get action by user")
	}
	if actor == nil {
		return nil, fluxgo.ErrorNotFound("Action by user not found")
	}

	user, err := h.userRepo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	if data.ActionBy != data.Id {
		return nil, utils.ErrorForbidden("You don't have permission to update this user", "user.forbidden")
	}

	validator := data.ValidateUpdateUserOptionalFields()

	fieldsToUpdate := 0
	repoData := &repositories.UpdateUserRepositoryReq{
		IdUser:   data.Id,
		ActionBy: data.ActionBy,
	}

	if validator.HasName && *data.Name != user.Name {
		repoData.Name = data.Name
		fieldsToUpdate++
	}

	if validator.HasPhoneNumber && *data.PhoneNumber != user.PhoneNumber {
		repoData.PhoneNumber = data.PhoneNumber
		fieldsToUpdate++
	}

	if validator.HasEmail && *data.Email != user.Email {
		userWithSameEmail, err := h.userRepo.GetUserByEmail(ctx, *data.Email)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get user by email")
		}
		if userWithSameEmail != nil && userWithSameEmail.IdUser != user.IdUser {
			return nil, fluxgo.ErrorBadRequest("Email already in use", "user.duplicate_email")
		}
		repoData.Email = data.Email
		fieldsToUpdate++
	}

	if validator.HasPassword {
		secret := utils.Secret{}
		hash, err := secret.Encrypt(h.config, *data.Password)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to encrypt password")
		}
		repoData.Password = &hash
		fieldsToUpdate++
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest(
			"At least one field must be different to update",
			"user.no_fields_to_update",
		)
	}

	err = h.userRepo.UpdateUser(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update user")
	}

	updatedUser, err := h.userRepo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get updated user")
	}
	if updatedUser == nil {
		return nil, fluxgo.ErrorNotFound("Updated user not found")
	}

	return &dto.UpdateUserRes{User: *updatedUser}, nil
}
