package handlers

import (
	c "context"
	"errors"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/user/dtos"
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerCreateUser struct {
	db             *fluxgo.Database
	config         *config.Env
	userRepo       *repositories.UserRepository
	addressRepo    *repositories.AddressRepository
	userBabyRepo   *repositories.UserBabyRepository
	consentLogRepo *repositories.ConsentLogRepository
}

type createUserTxError struct {
	err *fluxgo.GlobalError
}

func (e *createUserTxError) Error() string {
	return e.err.Message
}

func HandlerCreateUserStart(
	db *fluxgo.Database,
	config *config.Env,
	userRepo *repositories.UserRepository,
	addressRepo *repositories.AddressRepository,
	userBabyRepo *repositories.UserBabyRepository,
	consentLogRepo *repositories.ConsentLogRepository,
) *HandlerCreateUser {
	return &HandlerCreateUser{
		db,
		config,
		userRepo,
		addressRepo,
		userBabyRepo,
		consentLogRepo,
	}
}

func (h *HandlerCreateUser) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateUserReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateUser) Execute(ctx c.Context, data *dto.CreateUserReq) (*dto.CreateUserRes, *fluxgo.GlobalError) {
	secret := utils.Secret{}

	hashPassword, err := secret.Encrypt(h.config, data.Password)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to encrypt password")
	}

	userWithSameCPf, err := h.userRepo.GetUserByCpf(ctx, data.Cpf)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user by cpf")
	}
	if userWithSameCPf != nil {
		return nil, fluxgo.ErrorBadRequest("User with same cpf already exists", "user.duplicate_cpf")
	}

	userWithSameEmail, err := h.userRepo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user by email")
	}
	if userWithSameEmail != nil {
		return nil, fluxgo.ErrorBadRequest("User with same email already exists", "user.duplicate_email")
	}

	userWithSamePhoneNumber, err := h.userRepo.GetUserByPhoneNumber(ctx, data.PhoneNumber)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user by phone number")
	}
	if userWithSamePhoneNumber != nil {
		return nil, fluxgo.ErrorBadRequest("User with same phone number already exists", "user.duplicate_phone_number")
	}

	idUser := utils.IdGenerate(utils.UserEntity)

	req := repositories.CreateUserRepositoryReq{
		IdUser:      idUser,
		Name:        data.Name,
		Cpf:         data.Cpf,
		Type:        data.Type,
		Email:       data.Email,
		Password:    hashPassword,
		PhoneNumber: data.PhoneNumber,
	}

	validator := data.ValidateCreateUserOptionalFields()

	switch data.Type {
	case entities.EnumUserTypeCommon:
		return h.handleCommon(ctx, data, &req, validator, idUser)
	case entities.EnumUserTypeAdmin, entities.EnumUserTypeNurse:
		return h.handleWorker(ctx, data, &req, validator, idUser)
	}

	return nil, fluxgo.ErrorBadRequest("Invalid user type", "user.invalid_type")
}

func (h *HandlerCreateUser) handleCommon(ctx c.Context, data *dto.CreateUserReq, req *repositories.CreateUserRepositoryReq, validator dto.CreateUserOptionalFields, idUser string) (*dto.CreateUserRes, *fluxgo.GlobalError) {
	if !validator.CanCreateCommon {
		return nil, fluxgo.ErrorBadRequest("Missing required fields for common user", "user.missing_fields")
	}

	if utils.IsFutureDate(*data.BirthDate) {
		return nil, fluxgo.ErrorBadRequest("Birth date cannot be in the future", "user.invalid_birth_date")
	}

	birthDateTime, err := utils.StringToDate(*data.BirthDate)
	if err != nil {
		return nil, fluxgo.ErrorBadRequest("Invalid birth date format", "user.invalid_birth_date_format")
	}

	req.BirthDate = *birthDateTime
	req.ActionBy = idUser

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.userRepo.CreateUserTx(ctx, tx, req)
		if err != nil {
			return fmt.Errorf("error to create user: %w", err)
		}

		errHandler := h.createAddress(ctx, &dto.CreateAddressReq{
			ActionBy:          idUser,
			AddressCreateBase: *data.Address,
		}, tx)
		if errHandler != nil {
			return &createUserTxError{err: errHandler}
		}

		errHandler = h.createConsentLog(ctx, &dto.CreateConsentReq{
			ActionBy:          idUser,
			CreateConsentBase: *data.ConsentLog,
		}, tx)
		if errHandler != nil {
			return &createUserTxError{err: errHandler}
		}

		if data.UserBaby != nil {
			for _, userBaby := range *data.UserBaby {
				errHandler := h.createUserBaby(ctx, &dto.CreateUserBabyReq{
					ActionBy:           idUser,
					UserBabyCreateBase: userBaby,
				}, tx)
				if errHandler != nil {
					return &createUserTxError{err: errHandler}
				}
			}
		}

		return nil
	})
	if err != nil {
		var txErr *createUserTxError
		if errors.As(err, &txErr) {
			return nil, txErr.err
		}
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	user, err := h.userRepo.GetUserById(ctx, idUser)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	return &dto.CreateUserRes{UserOut: sharedDto.NewUserOut(*user)}, nil
}

func (h *HandlerCreateUser) handleWorker(ctx c.Context, data *dto.CreateUserReq, req *repositories.CreateUserRepositoryReq, validator dto.CreateUserOptionalFields, idUser string) (*dto.CreateUserRes, *fluxgo.GlobalError) {
	if !validator.CanCreateWorker {
		return nil, fluxgo.ErrorBadRequest("Missing required fields for worker user", "user.missing_fields")
	}

	user, err := h.userRepo.GetUserById(ctx, *data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden("User does not have permission to create another user", "user.forbidden")
	}

	req.ActionBy = *data.ActionBy
	req.InternalIdentifier = data.Identifier

	err = h.userRepo.CreateUser(ctx, req)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create user")
	}

	user, err = h.userRepo.GetUserById(ctx, idUser)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	return &dto.CreateUserRes{UserOut: sharedDto.NewUserOut(*user)}, nil
}

func (h *HandlerCreateUser) createConsentLog(ctx c.Context, data *dto.CreateConsentReq, tx *sqlx.Tx) *fluxgo.GlobalError {
	idConsentLog := utils.IdGenerate(utils.ConsentLogEntity)

	repoData := &repositories.CreateConsentRepositoryReq{
		TermsVersion: data.TermsVersion,
		Ip:           data.IpAddress,
		UserAgent:    data.UserAgent,
		IdUser:       data.ActionBy,
		IdConsentLog: idConsentLog,
	}

	err := h.consentLogRepo.CreateConsentLogTx(ctx, tx, repoData)
	if err != nil {
		return fluxgo.ErrorInternalError("Error to create consent")
	}

	return nil
}

func (h *HandlerCreateUser) createAddress(ctx c.Context, data *dto.CreateAddressReq, tx *sqlx.Tx) *fluxgo.GlobalError {
	addressData, err := utils.GetAddressByZipCodeOptionalCoordinates(ctx, data.ZipCode, h.config)
	if err != nil {
		return fluxgo.ErrorInternalError(err.Error())
	}

	idAddress := utils.IdGenerate(utils.AddressEntity)

	repoData := &repositories.CreateAddressRepositoryReq{
		IdAddress:    idAddress,
		IdUser:       &data.ActionBy,
		Zipcode:      data.ZipCode,
		Street:       addressData.Street,
		Number:       data.Number,
		City:         addressData.City,
		State:        addressData.State,
		Neighborhood: addressData.Neighborhood,
		Complement:   data.Complement,
		Latitude:     addressData.Latitude,
		Longitude:    addressData.Longitude,
	}

	err = h.addressRepo.CreateAddressTx(ctx, tx, repoData)
	if err != nil {
		return fluxgo.ErrorInternalError("Error to create address")
	}

	return nil
}

func (h *HandlerCreateUser) createUserBaby(ctx c.Context, data *dto.CreateUserBabyReq, tx *sqlx.Tx) *fluxgo.GlobalError {
	if utils.IsFutureDate(data.BirthDate) {
		return fluxgo.ErrorBadRequest("Birth date cannot be in the future", "user_baby.invalid_birth_date")
	}

	birthDateTime, err := utils.StringToDate(data.BirthDate)
	if err != nil {
		return fluxgo.ErrorBadRequest("Invalid birth date format", "user_baby.invalid_birth_date_format")
	}

	userBabyId := utils.IdGenerate(utils.UserBabyEntity)

	err = h.userBabyRepo.CreateUserBabyTx(ctx, tx, &repositories.CreateUserBabyRepositoryReq{
		IdUserBaby: userBabyId,
		IdUser:     data.ActionBy,
		Name:       data.Name,
		BirthDate:  *birthDateTime,
	})
	if err != nil {
		return fluxgo.ErrorInternalError("Error to create baby")
	}

	return nil
}
