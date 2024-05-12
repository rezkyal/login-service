// This file contains the interfaces for the usecase layer.
// The usecase layer is responsible for handling business logic.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package usecase

import (
	"context"
)

type UsecaseInterface interface {
	RegisterNewUser(ctx context.Context, input RegisterNewUserInput) (RegisterNewUserOutput, error)
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)
	GetUserData(ctx context.Context, input GetUserDataInput) (GetUserDataOutput, error)
	UpdateUserData(ctx context.Context, input UpdateUserDataInput) (UpdateUserDataOutput, error)
}
