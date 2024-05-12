// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	InsertNewUser(ctx context.Context, input InsertNewUserInput) (InsertNewUserOutput, error)
	GetPasswordByPhoneNumber(ctx context.Context, input GetPasswordByPhoneNumberInput) (output GetPasswordByPhoneNumberOutput, err error)
	GetUserDataById(ctx context.Context, input GetUserDataByIdInput) (output GetUserDataByIdOutput, err error)
	UpdateUserData(ctx context.Context, input UpdateUserDataInput) (UpdateUserDataOutput, error)
	UpdateTotalLoginById(ctx context.Context, input UpdateTotalLoginByIdInput) (err error)
}
