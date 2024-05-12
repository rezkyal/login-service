package repository

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func TestRepository_InsertNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		input InsertNewUserInput
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		want     InsertNewUserOutput
		wantErr  bool
	}{
		{
			name: "Success, confict phone number exists",
			args: args{
				input: InsertNewUserInput{
					PhoneNumber: "12345",
					FullName:    "fullname",
					Password:    "password",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(InsertNewUserQuery)).
					WithArgs(a.input.PhoneNumber, a.input.FullName, a.input.Password).
					WillReturnError(&pq.Error{
						Code: "23505",
					})
			},
			want: InsertNewUserOutput{
				IsPhoneNumberExists: true,
			},
			wantErr: false,
		},
		{
			name: "Error when query",
			args: args{
				input: InsertNewUserInput{
					PhoneNumber: "12345",
					FullName:    "fullname",
					Password:    "password",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(InsertNewUserQuery)).
					WithArgs(a.input.PhoneNumber, a.input.FullName, a.input.Password).
					WillReturnError(errors.New("test"))
			},
			want:    InsertNewUserOutput{},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				input: InsertNewUserInput{
					PhoneNumber: "12345",
					FullName:    "fullname",
					Password:    "password",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(InsertNewUserQuery)).
					WithArgs(a.input.PhoneNumber, a.input.FullName, a.input.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(int64(50)))
			},
			want: InsertNewUserOutput{
				Id: 50,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			r := &Repository{
				Db: db,
			}
			got, err := r.InsertNewUser(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.InsertNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.InsertNewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_UpdateUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

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
			name: "Success, conflict phone number exists",
			args: args{
				input: UpdateUserDataInput{
					Id:          10,
					PhoneNumber: "phone_number",
					FullName:    "full_name",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(UpdateUserDataQuery)).
					WithArgs(a.input.Id, a.input.PhoneNumber, a.input.FullName, a.input.Id).
					WillReturnError(&pq.Error{
						Code: "23505",
					})
			},
			want: UpdateUserDataOutput{
				IsPhoneNumberExists: true,
			},
			wantErr: false,
		},
		{
			name: "Error when query",
			args: args{
				input: UpdateUserDataInput{
					Id:          10,
					PhoneNumber: "phone_number",
					FullName:    "full_name",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(UpdateUserDataQuery)).
					WithArgs(a.input.Id, a.input.PhoneNumber, a.input.FullName, a.input.Id).
					WillReturnError(errors.New("test"))
			},
			want:    UpdateUserDataOutput{},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				input: UpdateUserDataInput{
					Id:          10,
					PhoneNumber: "phone_number",
					FullName:    "full_name",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(UpdateUserDataQuery)).
					WithArgs(a.input.Id, a.input.PhoneNumber, a.input.FullName, a.input.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    UpdateUserDataOutput{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			r := &Repository{
				Db: db,
			}
			got, err := r.UpdateUserData(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.UpdateUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetPasswordByPhoneNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		input GetPasswordByPhoneNumberInput
	}
	tests := []struct {
		name       string
		args       args
		mockFunc   func(args)
		wantOutput GetPasswordByPhoneNumberOutput
		wantErr    bool
	}{
		{
			name: "Error when query",
			args: args{
				input: GetPasswordByPhoneNumberInput{
					PhoneNumber: "12345",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPasswordByPhoneNumberQuery)).
					WithArgs(a.input.PhoneNumber).
					WillReturnError(errors.New("test"))
			},
			wantOutput: GetPasswordByPhoneNumberOutput{},
			wantErr:    true,
		},
		{
			name: "Success",
			args: args{
				input: GetPasswordByPhoneNumberInput{
					PhoneNumber: "12345",
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPasswordByPhoneNumberQuery)).
					WithArgs(a.input.PhoneNumber).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password", "phone_number"}).
						AddRow(int64(50), "passwordaa", "phone_000"))
			},
			wantOutput: GetPasswordByPhoneNumberOutput{
				Id:          50,
				PhoneNumber: "phone_000",
				Password:    "passwordaa",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			r := &Repository{
				Db: db,
			}
			gotOutput, err := r.GetPasswordByPhoneNumber(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetPasswordByPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.GetPasswordByPhoneNumber() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_UpdateTotalLoginById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		input UpdateTotalLoginByIdInput
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(args)
		wantErr  bool
	}{
		{
			name: "Error when query",
			args: args{
				input: UpdateTotalLoginByIdInput{
					Id: 99,
				},
			},
			mockFunc: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(UpdateTotalLoginById)).
					WithArgs(a.input.Id).
					WillReturnError(errors.New("test"))
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				input: UpdateTotalLoginByIdInput{
					Id: 99,
				},
			},
			mockFunc: func(a args) {
				mock.ExpectExec(regexp.QuoteMeta(UpdateTotalLoginById)).
					WithArgs(a.input.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.mockFunc(tt.args)
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: db,
			}
			if err := r.UpdateTotalLoginById(context.Background(), tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateTotalLoginById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetUserDataById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		input GetUserDataByIdInput
	}
	tests := []struct {
		name       string
		args       args
		mockFunc   func(args)
		wantOutput GetUserDataByIdOutput
		wantErr    bool
	}{
		{
			name: "Error when query",
			args: args{
				input: GetUserDataByIdInput{
					Id: 34,
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(GetUserDataByIdQuery)).
					WithArgs(a.input.Id).
					WillReturnError(errors.New("test"))
			},
			wantOutput: GetUserDataByIdOutput{},
			wantErr:    true,
		},
		{
			name: "Success",
			args: args{
				input: GetUserDataByIdInput{
					Id: 34,
				},
			},
			mockFunc: func(a args) {
				mock.ExpectQuery(regexp.QuoteMeta(GetUserDataByIdQuery)).
					WithArgs(a.input.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "phone_number"}).
						AddRow("50", "fullname", "phone_000"))
			},
			wantOutput: GetUserDataByIdOutput{
				Id:          "50",
				FullName:    "fullname",
				PhoneNumber: "phone_000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)
			r := &Repository{
				Db: db,
			}
			gotOutput, err := r.GetUserDataById(context.Background(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserDataById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.GetUserDataById() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
