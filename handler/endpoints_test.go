package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServer_Registration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := usecase.NewMockUsecaseInterface(ctrl)

	type args struct {
		ctx func() (echo.Context, *httptest.ResponseRecorder)
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		respFunc func(*httptest.ResponseRecorder) interface{}
		wantCode int
		wantResp interface{}
		wantErr  bool
	}{
		{
			name: "error validations",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "08123456")
					data.Set("full_name", "fu")
					data.Set("password", "asfa")

					req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.ValidationErrorsResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				sort.Slice(resp, func(i, j int) bool {
					return resp[i].Field < resp[j].Field
				})

				return resp
			},
			wantCode: http.StatusBadRequest,
			wantResp: generated.ValidationErrorsResponse{
				generated.ValidationError{
					Field:   "full_name",
					Message: "must be at minimum 3 characters and maximum 60 characters",
				},
				generated.ValidationError{
					Field:   "password",
					Message: "must be at minimum 6 characters and maximum 64 characters & must containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters",
				},
				generated.ValidationError{
					Field:   "phone_number",
					Message: "must be at minimum 10 characters and maximum 13 characters & must start with the Indonesia country code “+62”",
				},
			},
			wantErr: false,
		},
		{
			name: "error when RegisterNewUser",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("full_name", "fullloooo")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().RegisterNewUser(gomock.Any(), gomock.Eq(
					usecase.RegisterNewUserInput{
						PhoneNumber: "+62812345678",
						FullName:    "fullloooo",
						Password:    "AAssff1!",
					},
				)).Return(usecase.RegisterNewUserOutput{}, errors.New("test"))
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusInternalServerError,
			wantResp: generated.BasicErrorResponse{
				Message: "Internal server error",
			},
			wantErr: false,
		},
		{
			name: "error phone number already exists",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("full_name", "fullloooo")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().RegisterNewUser(gomock.Any(), gomock.Eq(
					usecase.RegisterNewUserInput{
						PhoneNumber: "+62812345678",
						FullName:    "fullloooo",
						Password:    "AAssff1!",
					},
				)).Return(usecase.RegisterNewUserOutput{
					IsPhoneNumberExists: true,
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusConflict,
			wantResp: generated.BasicErrorResponse{
				Message: "Phone number already used",
			},
			wantErr: false,
		},
		{
			name: "successful",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("full_name", "fullloooo")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().RegisterNewUser(gomock.Any(), gomock.Eq(
					usecase.RegisterNewUserInput{
						PhoneNumber: "+62812345678",
						FullName:    "fullloooo",
						Password:    "AAssff1!",
					},
				)).Return(usecase.RegisterNewUserOutput{
					Id: 5,
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.SuccessRegistrationResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusOK,
			wantResp: generated.SuccessRegistrationResponse{
				Id: "5",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			s := NewServer(NewServerOptions{
				Usecase: mockUsecase,
			})

			ctx, rec := tt.args.ctx()

			if err := s.Registration(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.Registration() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantCode, rec.Code)

			resp := tt.respFunc(rec)

			if respRaw, ok := resp.(generated.ValidationErrorsResponse); ok {
				sort.Slice(respRaw, func(i, j int) bool {
					return respRaw[i].Field < respRaw[j].Field
				})
			}

			assert.Equal(t, tt.wantResp, resp)
		})
	}
}

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := usecase.NewMockUsecaseInterface(ctrl)

	type args struct {
		ctx func() (echo.Context, *httptest.ResponseRecorder)
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		respFunc func(*httptest.ResponseRecorder) interface{}
		wantCode int
		wantResp interface{}
		wantErr  bool
	}{
		{
			name: "error when Login",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Eq(usecase.LoginInput{
					PhoneNumber: "+62812345678",
					Password:    "AAssff1!",
				})).Return(usecase.LoginOutput{}, errors.New("test"))
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusInternalServerError,
			wantResp: generated.BasicErrorResponse{
				Message: "Internal server error",
			},
			wantErr: false,
		},
		{
			name: "error data not found",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Eq(usecase.LoginInput{
					PhoneNumber: "+62812345678",
					Password:    "AAssff1!",
				})).Return(usecase.LoginOutput{
					IsDataNotFound: true,
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusBadRequest,
			wantResp: generated.BasicErrorResponse{
				Message: "Phone number not found",
			},
			wantErr: false,
		},
		{
			name: "error password wrong",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Eq(usecase.LoginInput{
					PhoneNumber: "+62812345678",
					Password:    "AAssff1!",
				})).Return(usecase.LoginOutput{
					IsPasswordWrong: true,
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusBadRequest,
			wantResp: generated.BasicErrorResponse{
				Message: "Wrong password",
			},
			wantErr: false,
		},
		{
			name: "success",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					data := url.Values{}
					data.Set("phone_number", "+62812345678")
					data.Set("password", "AAssff1!")

					req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Eq(usecase.LoginInput{
					PhoneNumber: "+62812345678",
					Password:    "AAssff1!",
				})).Return(usecase.LoginOutput{
					Token: "tokennn",
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.LoginSuccessResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusOK,
			wantResp: generated.LoginSuccessResponse{
				Message: "Login success",
				Token:   "tokennn",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			s := NewServer(NewServerOptions{
				Usecase: mockUsecase,
			})

			ctx, rec := tt.args.ctx()

			if err := s.Login(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.Login() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantCode, rec.Code)

			resp := tt.respFunc(rec)
			assert.Equal(t, tt.wantResp, resp)
		})
	}
}

func TestServer_ProfileGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := usecase.NewMockUsecaseInterface(ctrl)

	utils.KeyDataPrivate, _ = os.ReadFile("./../rsakey/jwtrsa256.key")
	utils.KeyDataPublic, _ = os.ReadFile("./../rsakey/jwtrsa256.key.pub")

	type args struct {
		ctx func() (echo.Context, *httptest.ResponseRecorder)
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		respFunc func(*httptest.ResponseRecorder) interface{}
		wantCode int
		wantResp interface{}
		wantErr  bool
	}{
		{
			name: "Error token invalid",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					// token, _ := utils.GenerateToken(50)

					token := "abcd"

					token = fmt.Sprintf("Bearer %s", token)

					req := httptest.NewRequest(http.MethodGet, "/profile", nil)
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusForbidden,
			wantResp: generated.BasicErrorResponse{
				Message: "Forbiddenn",
			},
			wantErr: false,
		},
		{
			name: "Error when GetUserData",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					token, _ := utils.GenerateToken(50)

					token = fmt.Sprintf("Bearer %s", token)

					req := httptest.NewRequest(http.MethodGet, "/profile", nil)
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().GetUserData(gomock.Any(), gomock.Eq(usecase.GetUserDataInput{
					Id: 50,
				})).Return(usecase.GetUserDataOutput{}, errors.New("test"))
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusInternalServerError,
			wantResp: generated.BasicErrorResponse{
				Message: "Internal server error",
			},
			wantErr: false,
		},
		{
			name: "Success",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					token, _ := utils.GenerateToken(50)

					token = fmt.Sprintf("Bearer %s", token)

					req := httptest.NewRequest(http.MethodGet, "/profile", nil)
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().GetUserData(gomock.Any(), gomock.Eq(usecase.GetUserDataInput{
					Id: 50,
				})).Return(usecase.GetUserDataOutput{
					PhoneNumber: "123456789",
					FullName:    "fullnamee",
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.ProfileGetResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusOK,
			wantResp: generated.ProfileGetResponse{
				PhoneNumber: "123456789",
				FullName:    "fullnamee",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			s := NewServer(NewServerOptions{
				Usecase: mockUsecase,
			})

			ctx, rec := tt.args.ctx()

			if err := s.ProfileGet(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.ProfileGet() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantCode, rec.Code)

			resp := tt.respFunc(rec)
			assert.Equal(t, tt.wantResp, resp)
		})
	}
}

func TestServer_ProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := usecase.NewMockUsecaseInterface(ctrl)

	utils.KeyDataPrivate, _ = os.ReadFile("./../rsakey/jwtrsa256.key")
	utils.KeyDataPublic, _ = os.ReadFile("./../rsakey/jwtrsa256.key.pub")

	type args struct {
		ctx func() (echo.Context, *httptest.ResponseRecorder)
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		respFunc func(*httptest.ResponseRecorder) interface{}
		wantCode int
		wantResp interface{}
		wantErr  bool
	}{
		{
			name: "Error when login",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					// token, _ := utils.GenerateToken(50)
					token := "abc"

					token = fmt.Sprintf("Bearer %s", token)

					req := httptest.NewRequest(http.MethodPost, "/profile", nil)
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusForbidden,
			wantResp: generated.BasicErrorResponse{
				Message: "Forbidden",
			},
			wantErr: false,
		},
		{
			name: "Error validations",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					token, _ := utils.GenerateToken(50)

					token = fmt.Sprintf("Bearer %s", token)

					data := url.Values{}
					data.Set("phone_number", "08123456")
					data.Set("full_name", "fu")

					req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.ValidationErrorsResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusBadRequest,
			wantResp: generated.ValidationErrorsResponse{
				generated.ValidationError{
					Field:   "full_name",
					Message: "must be at minimum 3 characters and maximum 60 characters",
				},
				generated.ValidationError{
					Field:   "phone_number",
					Message: "must be at minimum 10 characters and maximum 13 characters & must start with the Indonesia country code “+62”",
				},
			},
			wantErr: false,
		},
		{
			name: "Error when UpdateUserData",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					token, _ := utils.GenerateToken(50)

					token = fmt.Sprintf("Bearer %s", token)

					data := url.Values{}
					data.Set("phone_number", "+628123456784")
					data.Set("full_name", "fullnameeaa")

					req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().UpdateUserData(gomock.Any(), gomock.Eq(usecase.UpdateUserDataInput{
					Id:          50,
					PhoneNumber: "+628123456784",
					FullName:    "fullnameeaa",
				})).Return(usecase.UpdateUserDataOutput{}, errors.New("test"))
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusInternalServerError,
			wantResp: generated.BasicErrorResponse{
				Message: "Internal server error",
			},
			wantErr: false,
		},
		{
			name: "Error phone number exists",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					token, _ := utils.GenerateToken(50)

					token = fmt.Sprintf("Bearer %s", token)

					data := url.Values{}
					data.Set("phone_number", "+628123456784")
					data.Set("full_name", "fullnameeaa")

					req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().UpdateUserData(gomock.Any(), gomock.Eq(usecase.UpdateUserDataInput{
					Id:          50,
					PhoneNumber: "+628123456784",
					FullName:    "fullnameeaa",
				})).Return(usecase.UpdateUserDataOutput{
					IsPhoneNumberExists: true,
				}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusConflict,
			wantResp: generated.BasicErrorResponse{
				Message: "Phone number already used",
			},
			wantErr: false,
		},
		{
			name: "Success",
			args: args{
				ctx: func() (echo.Context, *httptest.ResponseRecorder) {
					e := echo.New()

					token, _ := utils.GenerateToken(50)

					token = fmt.Sprintf("Bearer %s", token)

					data := url.Values{}
					data.Set("phone_number", "+628123456784")
					data.Set("full_name", "fullnameeaa")

					req := httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(data.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Add("Authorization", token)
					rec := httptest.NewRecorder()

					return e.NewContext(req, rec), rec
				},
			},
			mockFunc: func(a args) {
				mockUsecase.EXPECT().UpdateUserData(gomock.Any(), gomock.Eq(usecase.UpdateUserDataInput{
					Id:          50,
					PhoneNumber: "+628123456784",
					FullName:    "fullnameeaa",
				})).Return(usecase.UpdateUserDataOutput{}, nil)
			},
			respFunc: func(rec *httptest.ResponseRecorder) interface{} {
				var resp generated.BasicSuccessResponse
				json.Unmarshal(rec.Body.Bytes(), &resp)

				return resp
			},
			wantCode: http.StatusOK,
			wantResp: generated.BasicSuccessResponse{
				Message: "Update success",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			s := NewServer(NewServerOptions{
				Usecase: mockUsecase,
			})

			ctx, rec := tt.args.ctx()
			if err := s.ProfileUpdate(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.ProfileUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantCode, rec.Code)

			resp := tt.respFunc(rec)

			if respRaw, ok := resp.(generated.ValidationErrorsResponse); ok {
				sort.Slice(respRaw, func(i, j int) bool {
					return respRaw[i].Field < respRaw[j].Field
				})
			}

			assert.Equal(t, tt.wantResp, resp)
		})
	}
}
