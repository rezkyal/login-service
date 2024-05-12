package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *Usecase) RegisterNewUser(ctx context.Context, input RegisterNewUserInput) (RegisterNewUserOutput, error) {
	var (
		hashCostStr = os.Getenv("BCRYPT_COST")
	)

	hashCost, err := strconv.Atoi(hashCostStr)
	if err != nil {
		log.Println("[WARN][RegisterNewUser] error when converting hashCost", err)
	}

	if hashCost == 0 {
		hashCost = 5
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), hashCost)
	if err != nil {
		return RegisterNewUserOutput{}, fmt.Errorf("[RegisterNewUser] error when hash password, err: %+v", err)
	}

	output, err := u.Repository.InsertNewUser(ctx, repository.InsertNewUserInput{
		PhoneNumber: input.PhoneNumber,
		FullName:    input.FullName,
		Password:    string(hashedPassword),
	})

	if err != nil {
		return RegisterNewUserOutput{}, fmt.Errorf("[RegisterNewUser] error when InsertNewUser %+v", err)
	}

	if output.IsPhoneNumberExists {
		return RegisterNewUserOutput{
			IsPhoneNumberExists: true,
		}, nil
	}

	return RegisterNewUserOutput{
		Id: output.Id,
	}, nil
}

func (u *Usecase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	passwordRes, err := u.Repository.GetPasswordByPhoneNumber(ctx, repository.GetPasswordByPhoneNumberInput{
		PhoneNumber: input.PhoneNumber,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return LoginOutput{
				IsDataNotFound: true,
			}, nil
		}

		return LoginOutput{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordRes.Password), []byte(input.Password))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return LoginOutput{
				IsPasswordWrong: true,
			}, nil
		}

		return LoginOutput{}, err
	}

	jwtToken, err := utils.GenerateToken(passwordRes.Id)

	if err != nil {
		return LoginOutput{}, err
	}

	go func(id int64) {
		err := u.Repository.UpdateTotalLoginById(context.Background(), repository.UpdateTotalLoginByIDInput{
			Id: id,
		})

		if err != nil {
			log.Println("[ERROR][Login] error when UpdateTotalLoginById", err)
		}
	}(passwordRes.Id)

	return LoginOutput{
		Token: jwtToken,
	}, nil
}

func (u *Usecase) GetUserData(ctx context.Context, input GetUserDataInput) (GetUserDataOutput, error) {
	outputRepo, err := u.Repository.GetUserDataById(ctx, repository.GetUserDataByIdInput{
		Id: input.Id,
	})

	if err != nil {
		return GetUserDataOutput{}, err
	}

	return GetUserDataOutput{
		PhoneNumber: outputRepo.PhoneNumber,
		FullName:    outputRepo.FullName,
	}, nil
}

func (u *Usecase) UpdateUserData(ctx context.Context, input UpdateUserDataInput) (UpdateUserDataOutput, error) {
	userData, err := u.Repository.GetUserDataById(ctx, repository.GetUserDataByIdInput{
		Id: input.Id,
	})

	if err != nil {
		return UpdateUserDataOutput{}, err
	}

	if input.PhoneNumber != "" {
		userData.PhoneNumber = input.PhoneNumber
	}

	if input.FullName != "" {
		userData.FullName = input.FullName
	}

	outputRepo, err := u.Repository.UpdateUserData(ctx, repository.UpdateUserDataInput{
		Id:          input.Id,
		PhoneNumber: input.PhoneNumber,
		FullName:    input.FullName,
	})

	if err != nil {
		return UpdateUserDataOutput{}, err
	}

	return UpdateUserDataOutput{
		Id:                  outputRepo.Id,
		IsPhoneNumberExists: outputRepo.IsPhoneNumberExists,
	}, nil
}
