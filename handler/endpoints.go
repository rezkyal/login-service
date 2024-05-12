package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
)

// This endpoint is used to register a new user
// (POST /registration)
func (s *Server) Registration(ctx echo.Context) error {
	var (
		req  generated.RegistrationFormdataBody
		errs = make(map[string]error)
	)

	ctx.Bind(&req)

	errValidation := utils.ValidateFullName(req.FullName)
	if errValidation != nil {
		errs[FULLNAME_FIELD] = errValidation
	}

	errValidation = utils.ValidatePassword(req.Password)
	if errValidation != nil {
		errs[PASSWORD_FIELD] = errValidation
	}

	errValidation = utils.ValidatePhoneNumbers(req.PhoneNumber)
	if errValidation != nil {
		errs[PHONE_NUMBER_FIELD] = errValidation
	}

	if len(errs) != 0 {
		var (
			resp generated.ValidationErrorsResponse
		)

		for k, v := range errs {
			resp = append(resp, generated.ValidationError{
				Field:   k,
				Message: v.Error(),
			})
		}

		return ctx.JSON(http.StatusBadRequest, resp)
	}

	resp, err := s.Usecase.RegisterNewUser(ctx.Request().Context(), usecase.RegisterNewUserInput{
		PhoneNumber: req.PhoneNumber,
		FullName:    req.FullName,
		Password:    req.Password,
	})

	if err != nil {
		log.Println("[ERROR][Registration] error when RegisterNewUser", err)
		return ctx.JSON(http.StatusInternalServerError, generated.BasicErrorResponse{
			Message: "Internal server error",
		})
	}

	if resp.IsPhoneNumberExists {
		return ctx.JSON(http.StatusConflict, generated.BasicErrorResponse{
			Message: "Phone number already used",
		})
	}

	return ctx.JSON(http.StatusOK, generated.SuccessRegistrationResponse{
		Id: strconv.FormatInt(resp.Id, 10),
	})
}

// This endpoint is used to login a user
// (POST /login)
func (s *Server) Login(ctx echo.Context) error {
	var (
		req generated.LoginFormdataBody
	)

	ctx.Bind(&req)

	resp, err := s.Usecase.Login(ctx.Request().Context(), usecase.LoginInput{
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})

	if err != nil {
		log.Println("[ERROR][Login] error when Login", err)
		return ctx.JSON(http.StatusInternalServerError, generated.BasicErrorResponse{
			Message: "Internal server error",
		})
	}

	if resp.IsDataNotFound {
		return ctx.JSON(http.StatusBadRequest, generated.BasicErrorResponse{
			Message: "Phone number not found",
		})
	}

	if resp.IsPasswordWrong {
		return ctx.JSON(http.StatusBadRequest, generated.BasicErrorResponse{
			Message: "Wrong password",
		})
	}

	return ctx.JSON(http.StatusOK, generated.LoginSuccessResponse{
		Message: "Login success",
		Token:   resp.Token,
	})
}

// Get profile data based on the jwt headers
// (GET /profile)
func (s *Server) ProfileGet(ctx echo.Context) error {

	id, err := utils.TokenValidity(ctx)

	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.BasicErrorResponse{
			Message: "Forbiddenn",
		})
	}

	userData, err := s.Usecase.GetUserData(ctx.Request().Context(), usecase.GetUserDataInput{
		Id: id,
	})

	if err != nil {
		log.Println("[ERROR][ProfileGet] error when GetUserData", err)
		return ctx.JSON(http.StatusInternalServerError, generated.BasicErrorResponse{
			Message: "Internal server error",
		})
	}

	return ctx.JSON(http.StatusOK, generated.ProfileGetResponse{
		PhoneNumber: userData.PhoneNumber,
		FullName:    userData.FullName,
	})
}

// Update profile data based on the request body and the jwt headers
// (PUT /profile)
func (s *Server) ProfileUpdate(ctx echo.Context) error {

	id, err := utils.TokenValidity(ctx)

	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.BasicErrorResponse{
			Message: "Forbidden",
		})
	}

	var (
		req  generated.ProfileUpdateFormdataBody
		errs = make(map[string]error)
	)

	ctx.Bind(&req)

	if req.FullName != "" {
		errValidation := utils.ValidateFullName(req.FullName)
		if errValidation != nil {
			errs[FULLNAME_FIELD] = errValidation
		}
	}

	if req.PhoneNumber != "" {
		errValidation := utils.ValidatePhoneNumbers(req.PhoneNumber)
		if errValidation != nil {
			errs[PHONE_NUMBER_FIELD] = errValidation
		}
	}

	if len(errs) != 0 {
		var (
			resp generated.ValidationErrorsResponse
		)

		for k, v := range errs {
			resp = append(resp, generated.ValidationError{
				Field:   k,
				Message: v.Error(),
			})
		}

		return ctx.JSON(http.StatusBadRequest, resp)
	}

	output, err := s.Usecase.UpdateUserData(ctx.Request().Context(), usecase.UpdateUserDataInput{
		Id:          id,
		PhoneNumber: req.PhoneNumber,
		FullName:    req.FullName,
	})

	if err != nil {
		log.Println("[ERROR][ProfileUpdate] error when UpdateUserData", err)
		return ctx.JSON(http.StatusInternalServerError, generated.BasicErrorResponse{
			Message: "Internal server error",
		})
	}

	if output.IsPhoneNumberExists {
		return ctx.JSON(http.StatusConflict, generated.BasicErrorResponse{
			Message: "Phone number already used",
		})
	}

	return ctx.JSON(http.StatusOK, generated.BasicSuccessResponse{
		Message: "Update success",
	})
}
