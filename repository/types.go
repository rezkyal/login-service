// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	FullName string
}

type InsertNewUserInput struct {
	PhoneNumber string
	FullName    string
	Password    string
}

type InsertNewUserOutput struct {
	Id                  int64
	IsPhoneNumberExists bool
}

type UpdateUserDataInput struct {
	Id          int64
	PhoneNumber string
	FullName    string
}

type UpdateUserDataOutput struct {
	IsPhoneNumberExists bool
}

type GetPasswordByPhoneNumberInput struct {
	PhoneNumber string
}

type GetPasswordByPhoneNumberOutput struct {
	Id          int64
	PhoneNumber string
	Password    string
}

type GetUserDataByIdInput struct {
	Id int64
}

type GetUserDataByIdOutput struct {
	Id          string
	FullName    string
	PhoneNumber string
}

type UpdateTotalLoginByIdInput struct {
	Id int64
}
