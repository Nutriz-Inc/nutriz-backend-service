package tests

import (
	"net/http"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/module"
	"nutriz-backend-service/shared/utils"
	"strings"
	"testing"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	fx, app := module.Module().GetTestApp(t)
	defer fx.RequireStart().RequireStop()

	publicEndpoint := "/public/user"
	internalEndpoint := "/internal/user"

	makeAddress := func(zipcode string) *dto.AddressCreateBase {
		return &dto.AddressCreateBase{
			ZipCode:    zipcode,
			Number:     utils.StringPtr("123"),
			Complement: utils.StringPtr("Apto 10"),
		}
	}

	makeCommonBody := func(name, cpf, email, phone, birthDate, zipcode string) dto.CreateUserReq {
		return dto.CreateUserReq{
			Type:        entities.EnumUserTypeCommon,
			Name:        name,
			Cpf:         cpf,
			Email:       email,
			Password:    "12345678",
			BirthDate:   utils.StringPtr(birthDate),
			PhoneNumber: phone,
			Address:     makeAddress(zipcode),
			ConsentLog:  &ConsentLog,
		}
	}

	makeWorkerBody := func(userType entities.EnumUserType, name, cpf, email, phone, birthDate, identifier string) dto.CreateUserReq {
		return dto.CreateUserReq{
			Type:        userType,
			Name:        name,
			Cpf:         cpf,
			Email:       email,
			Password:    "12345678",
			PhoneNumber: phone,
			Identifier:  utils.StringPtr(identifier),
		}
	}

	t.Run("Success", func(t *testing.T) {
		t.Run("Create common user without baby", func(t *testing.T) {
			data := makeCommonBody(
				"Joana Souza",
				"57694673800",
				"joana.souza@email.com",
				"+5511991111111",
				"1990-01-01",
				"01001000",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusCreated, status)
			idUser, ok := body["id_user"].(string)
			assert.True(t, ok)
			assert.True(t, strings.HasPrefix(idUser, "usr_"))
			assert.Equal(t, "common", body["type"])
			assert.Equal(t, data.Name, body["name"])
			assert.Equal(t, data.Email, body["email"])
			assert.Equal(t, idUser, body["created_by"])
			assert.Nil(t, body["password"])
		})

		t.Run("Create common user with baby", func(t *testing.T) {
			data := makeCommonBody(
				"Ana Ribeiro",
				"66684931490",
				"ana.ribeiro@email.com",
				"+5511992222222",
				"1988-03-10",
				"01002000",
			)
			babyName := "Lia"
			data.UserBaby = &[]dto.UserBabyCreateBase{
				{
					Name:      &babyName,
					BirthDate: "2024-01-01",
				},
			}

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotNil(t, body["id_user"])
			assert.Equal(t, "common", body["type"])
			assert.Equal(t, data.Email, body["email"])
		})

		t.Run("Create common user with multiple babies", func(t *testing.T) {
			data := makeCommonBody(
				"Marcia Lima",
				"34111259693",
				"marcia.lima@email.com",
				"+5511992222233",
				"1985-05-15",
				"01002100",
			)
			firstBabyName := "Noah"
			secondBabyName := "Alice"
			data.UserBaby = &[]dto.UserBabyCreateBase{
				{
					Name:      &firstBabyName,
					BirthDate: "2023-01-01",
				},
				{
					Name:      &secondBabyName,
					BirthDate: "2024-06-15",
				},
			}

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotNil(t, body["id_user"])
			assert.Equal(t, "common", body["type"])
			assert.Equal(t, data.Email, body["email"])
		})

		t.Run("Create worker user", func(t *testing.T) {
			adminHeaders := &utils.TestHeadersAdmin
			data := makeWorkerBody(
				entities.EnumUserTypeNurse,
				"Camila Santos",
				"49657104017",
				"camila.santos@email.com",
				"+5511993333333",
				"1992-06-20",
				"WORKER-001",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", internalEndpoint, data, adminHeaders)

			assert.Equal(t, http.StatusCreated, status)
			assert.NotNil(t, body["id_user"])
			assert.Equal(t, "nurse", body["type"])
			assert.Equal(t, "WORKER-001", body["internal_identifier"])
			assert.Equal(t, "usr_2veL1FPpuXxUaZcFaEC57BfpxWS", body["created_by"])
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("User with same cpf already exists", func(t *testing.T) {
			data := makeCommonBody(
				"Duplicado CPF",
				"47046117012",
				"cpf.duplicado@email.com",
				"+5511994444444",
				"1991-02-15",
				"01003000",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "User with same cpf already exists", body["message"])
			assert.Equal(t, "user.duplicate_cpf", body["code"])
		})

		t.Run("User with same email already exists", func(t *testing.T) {
			data := makeCommonBody(
				"Duplicado Email",
				"08425431808",
				"maria@email.com",
				"+5511995555555",
				"1993-11-05",
				"01004000",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "User with same email already exists", body["message"])
			assert.Equal(t, "user.duplicate_email", body["code"])
		})

		t.Run("User with same phone number already exists", func(t *testing.T) {
			data := makeCommonBody(
				"Duplicado Telefone",
				"34470578070",
				"testeeeeeeeee@email.com",
				"+5511991111111",
				"1993-11-05",
				"01004000",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "User with same phone number already exists", body["message"])
			assert.Equal(t, "user.duplicate_phone_number", body["code"])
		})

		t.Run("Missing required fields for common user", func(t *testing.T) {
			data := makeCommonBody(
				"Sem Endereco",
				"67521647114",
				"sem.endereco@email.com",
				"+5511996666666",
				"1994-07-22",
				"01005000",
			)
			data.Address = nil

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Missing required fields for common user", body["message"])
			assert.Equal(t, "user.missing_fields", body["code"])
		})

		t.Run("Missing required fields for worker user", func(t *testing.T) {
			adminHeaders := &utils.TestHeadersAdmin
			data := makeWorkerBody(
				entities.EnumUserTypeAdmin,
				"Sem Identificador",
				"06892078591",
				"sem.identificador@email.com",
				"+5511997777777",
				"1992-09-10",
				"WORKER-002",
			)
			data.Identifier = nil

			status, body := fluxgo.RunTestRequest(app, "POST", internalEndpoint, data, adminHeaders)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Missing required fields for worker user", body["message"])
			assert.Equal(t, "user.missing_fields", body["code"])
		})

		t.Run("Birth date cannot be in the future", func(t *testing.T) {
			data := makeCommonBody(
				"Nascimento Futuro",
				"85116405419",
				"nascimento.futuro@email.com",
				"+5511998888888",
				utils.GetTodayFormattedDate(1),
				"01006000",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Birth date cannot be in the future", body["message"])
			assert.Equal(t, "user.invalid_birth_date", body["code"])
		})

		t.Run("User does not have permission to create another user", func(t *testing.T) {
			commonHeaders := &utils.TestHeaders
			data := makeWorkerBody(
				entities.EnumUserTypeNurse,
				"Sem Permissao",
				"31584595582",
				"sem.permissao@email.com",
				"+5511999999998",
				"1991-04-18",
				"WORKER-003",
			)

			status, body := fluxgo.RunTestRequest(app, "POST", internalEndpoint, data, commonHeaders)

			assert.Equal(t, http.StatusForbidden, status)
			assert.Equal(t, "User does not have permission to create another user", body["message"])
			assert.Equal(t, "user.forbidden", body["code"])
		})

		t.Run("user_baby.invalid_birth_date", func(t *testing.T) {
			data := makeCommonBody(
				"Bebe com data invalida",
				"08425431808",
				"bebe.data.invalida@email.com",
				"+5511999999999",
				"1995-02-12",
				"01007000",
			)
			babyName := "Baby Futuro"
			data.UserBaby = &[]dto.UserBabyCreateBase{
				{
					Name:      &babyName,
					BirthDate: utils.GetTodayFormattedDate(1),
				},
			}

			status, body := fluxgo.RunTestRequest(app, "POST", publicEndpoint, data, nil)

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, "Birth date cannot be in the future", body["message"])
			assert.Equal(t, "user_baby.invalid_birth_date", body["code"])
		})
	})
}
