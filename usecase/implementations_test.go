package usecase

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_RegisterNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	type args struct {
		input RegisterNewUserInput
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		want     RegisterNewUserOutput
		wantErr  bool
	}{
		{
			name: "error, when GenerateFromPassword",
			args: args{
				input: RegisterNewUserInput{
					PhoneNumber: "phone-000",
					FullName:    "fullname",
					Password:    "asfasfascojscoascjasockmasockasmofiajsfoasijfasokmcoaskcjasoifjasofkmasocijasofiajsokmascoaskjcoasifjas",
				},
			},
			mockFunc: func(a args) {},
			want:     RegisterNewUserOutput{},
			wantErr:  true,
		},
		{
			name: "error, when InsertNewUser",
			args: args{
				input: RegisterNewUserInput{
					PhoneNumber: "phone-000",
					FullName:    "fullname",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().InsertNewUser(gomock.Any(), gomock.Any()).Return(repository.InsertNewUserOutput{}, errors.New("test"))
			},
			want:    RegisterNewUserOutput{},
			wantErr: true,
		},
		{
			name: "success, phone number exists",
			args: args{
				input: RegisterNewUserInput{
					PhoneNumber: "phone-000",
					FullName:    "fullname",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().InsertNewUser(gomock.Any(), gomock.Any()).Return(repository.InsertNewUserOutput{
					IsPhoneNumberExists: true,
				}, nil)
			},
			want: RegisterNewUserOutput{
				IsPhoneNumberExists: true,
			},
			wantErr: false,
		},
		{
			name: "success",
			args: args{
				input: RegisterNewUserInput{
					PhoneNumber: "phone-000",
					FullName:    "fullname",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().InsertNewUser(gomock.Any(), gomock.Any()).Return(repository.InsertNewUserOutput{
					Id: 10,
				}, nil)
			},
			want: RegisterNewUserOutput{
				Id: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			u := NewUsecase(NewUsecaseOptions{
				Repository: mockRepository,
			})
			got, err := u.RegisterNewUser(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.RegisterNewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	utils.KeyDataPrivate, _ = os.ReadFile("./../rsakey/jwtrsa256.key")
	utils.KeyDataPublic, _ = os.ReadFile("./../rsakey/jwtrsa256.key.pub")

	type args struct {
		input LoginInput
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		want     LoginOutput
		wantId   int64
		wantErr  bool
	}{
		{
			name: "error when GetPasswordByPhoneNumber data not found",
			args: args{
				input: LoginInput{
					PhoneNumber: "phone",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetPasswordByPhoneNumber(gomock.Any(), gomock.Eq(repository.GetPasswordByPhoneNumberInput{
					PhoneNumber: a.input.PhoneNumber,
				})).Return(repository.GetPasswordByPhoneNumberOutput{}, sql.ErrNoRows)
			},
			want: LoginOutput{
				IsDataNotFound: true,
			},
			wantErr: false,
		},
		{
			name: "error when GetPasswordByPhoneNumber",
			args: args{
				input: LoginInput{
					PhoneNumber: "phone",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetPasswordByPhoneNumber(gomock.Any(), gomock.Eq(repository.GetPasswordByPhoneNumberInput{
					PhoneNumber: a.input.PhoneNumber,
				})).Return(repository.GetPasswordByPhoneNumberOutput{}, errors.New("test"))
			},
			want:    LoginOutput{},
			wantErr: true,
		},
		{
			name: "success, password mismatch",
			args: args{
				input: LoginInput{
					PhoneNumber: "phone",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetPasswordByPhoneNumber(gomock.Any(), gomock.Eq(repository.GetPasswordByPhoneNumberInput{
					PhoneNumber: a.input.PhoneNumber,
				})).Return(repository.GetPasswordByPhoneNumberOutput{
					Id:          10,
					PhoneNumber: "phone",
					Password:    "$2a$05$N9yncSBoAMWxz/nyW7APGuzRkXXGh27574xz2pF8dj4vm.In9T0SW",
				}, nil)
			},
			want: LoginOutput{
				IsPasswordWrong: true,
			},
			wantErr: false,
		},
		{
			name: "error when CompareHashAndPassword",
			args: args{
				input: LoginInput{
					PhoneNumber: "phone",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetPasswordByPhoneNumber(gomock.Any(), gomock.Eq(repository.GetPasswordByPhoneNumberInput{
					PhoneNumber: a.input.PhoneNumber,
				})).Return(repository.GetPasswordByPhoneNumberOutput{
					Id:          10,
					PhoneNumber: "phone",
					Password:    "abcd",
				}, nil)
			},
			want:    LoginOutput{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				input: LoginInput{
					PhoneNumber: "phone",
					Password:    "aaaa",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetPasswordByPhoneNumber(gomock.Any(), gomock.Eq(repository.GetPasswordByPhoneNumberInput{
					PhoneNumber: a.input.PhoneNumber,
				})).Return(repository.GetPasswordByPhoneNumberOutput{
					Id:          10,
					PhoneNumber: "phone",
					Password:    "$2a$05$WgWdo896B1Qc3VQRIm78X.rdwOFwEo7dB.bgIbAx8wOBNZCx1eJ2q",
				}, nil)

				mockRepository.EXPECT().UpdateTotalLoginById(gomock.Any(), gomock.Eq(repository.UpdateTotalLoginByIdInput{
					Id: 10,
				})).Return(errors.New("test")).AnyTimes()
			},
			want:    LoginOutput{},
			wantId:  10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			u := NewUsecase(NewUsecaseOptions{
				Repository: mockRepository,
			})
			got, err := u.Login(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			token := got.Token
			got.Token = ""

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Login() = %v, want %v", got, tt.want)
			}

			if token != "" {
				parse, _ := utils.TokenParse(token)
				assert.Equal(t, tt.wantId, parse)
			}
		})
	}
}

func TestUsecase_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	type args struct {
		input GetUserDataInput
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		want     GetUserDataOutput
		wantErr  bool
	}{
		{
			name: "error when GetUserDataById",
			args: args{
				input: GetUserDataInput{
					Id: 11,
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetUserDataById(gomock.Any(), gomock.Eq(repository.GetUserDataByIdInput{
					Id: a.input.Id,
				})).Return(repository.GetUserDataByIdOutput{}, errors.New("test"))
			},
			want:    GetUserDataOutput{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				input: GetUserDataInput{
					Id: 11,
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetUserDataById(gomock.Any(), gomock.Eq(repository.GetUserDataByIdInput{
					Id: a.input.Id,
				})).Return(repository.GetUserDataByIdOutput{
					PhoneNumber: "phoneNumber",
					FullName:    "fullName",
				}, nil)
			},
			want: GetUserDataOutput{
				PhoneNumber: "phoneNumber",
				FullName:    "fullName",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			u := NewUsecase(NewUsecaseOptions{
				Repository: mockRepository,
			})
			got, err := u.GetUserData(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.GetUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_UpdateUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	type args struct {
		input UpdateUserDataInput
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		want     UpdateUserDataOutput
		wantErr  bool
	}{
		{
			name: "error when GetUserDataById",
			args: args{
				input: UpdateUserDataInput{
					Id:          10,
					PhoneNumber: "numberPhone",
					FullName:    "fullname",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetUserDataById(gomock.Any(), gomock.Eq(repository.GetUserDataByIdInput{
					Id: a.input.Id,
				})).Return(repository.GetUserDataByIdOutput{}, errors.New("test"))
			},
			want:    UpdateUserDataOutput{},
			wantErr: true,
		},
		{
			name: "error when UpdateUserData",
			args: args{
				input: UpdateUserDataInput{
					Id:          10,
					PhoneNumber: "numberPhone",
					FullName:    "fullname",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetUserDataById(gomock.Any(), gomock.Eq(repository.GetUserDataByIdInput{
					Id: a.input.Id,
				})).Return(repository.GetUserDataByIdOutput{
					Id:          "10",
					PhoneNumber: "phoneNumber",
					FullName:    "nameFull",
				}, nil)

				mockRepository.EXPECT().UpdateUserData(gomock.Any(), gomock.Eq(repository.UpdateUserDataInput{
					Id:          a.input.Id,
					PhoneNumber: a.input.PhoneNumber,
					FullName:    a.input.FullName,
				})).Return(repository.UpdateUserDataOutput{
					IsPhoneNumberExists: true,
				}, errors.New("test"))
			},
			want:    UpdateUserDataOutput{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				input: UpdateUserDataInput{
					Id:          10,
					PhoneNumber: "numberPhone",
					FullName:    "fullname",
				},
			},
			mockFunc: func(a args) {
				mockRepository.EXPECT().GetUserDataById(gomock.Any(), gomock.Eq(repository.GetUserDataByIdInput{
					Id: a.input.Id,
				})).Return(repository.GetUserDataByIdOutput{
					Id:          "10",
					PhoneNumber: "phoneNumber",
					FullName:    "nameFull",
				}, nil)

				mockRepository.EXPECT().UpdateUserData(gomock.Any(), gomock.Eq(repository.UpdateUserDataInput{
					Id:          a.input.Id,
					PhoneNumber: a.input.PhoneNumber,
					FullName:    a.input.FullName,
				})).Return(repository.UpdateUserDataOutput{
					IsPhoneNumberExists: true,
				}, nil)
			},
			want: UpdateUserDataOutput{
				IsPhoneNumberExists: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			u := NewUsecase(NewUsecaseOptions{
				Repository: mockRepository,
			})
			got, err := u.UpdateUserData(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.UpdateUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.UpdateUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}
