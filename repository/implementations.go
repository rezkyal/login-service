package repository

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT full_name FROM test WHERE id = $1", input.Id).Scan(&output.FullName)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (r *Repository) InsertNewUser(ctx context.Context, input InsertNewUserInput) (InsertNewUserOutput, error) {
	var (
		output InsertNewUserOutput
		newId  int64
	)

	err := r.Db.QueryRowContext(ctx, "INSERT INTO users(phone_number, full_name, password) values ($1, $2, $3) returning id", input.PhoneNumber, input.FullName, input.Password).Scan(&newId)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok && pgerr.Code == "23505" {
			return InsertNewUserOutput{
				IsPhoneNumberExists: true,
			}, nil
		}
		return output, errors.WithStack(err)
	}

	output.Id = newId

	return output, nil
}

func (r *Repository) UpdateUserData(ctx context.Context, input UpdateUserDataInput) (UpdateUserDataOutput, error) {
	_, err := r.Db.ExecContext(ctx, `UPDATE users 
		set phone_number = $2, 
		full_name = $3
		WHERE id = $1`, input.Id, input.PhoneNumber, input.FullName)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok && pgerr.Code == KEY_CONFLICT {
			return UpdateUserDataOutput{
				IsPhoneNumberExists: true,
			}, nil
		}
	}
	return UpdateUserDataOutput{}, errors.WithStack(err)
}

func (r *Repository) GetPasswordByPhoneNumber(ctx context.Context, input GetPasswordByPhoneNumberInput) (output GetPasswordByPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, password, phone_number FROM users WHERE phone_number = $1", input.PhoneNumber).Scan(&output.Id, &output.Password, &output.PhoneNumber)
	err = errors.WithStack(err)
	return
}

func (r *Repository) UpdateTotalLoginById(ctx context.Context, input UpdateTotalLoginByIDInput) (err error) {
	_, err = r.Db.ExecContext(ctx, `UPDATE users
		SET total_login = total_login + 1
		WHERE id = $1`, input.Id)

	errors.WithStack(err)
	return err
}

func (r *Repository) GetUserDataById(ctx context.Context, input GetUserDataByIdInput) (output GetUserDataByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number FROM users WHERE id = $1", input.Id).Scan(&output.Id, &output.FullName, &output.PhoneNumber)
	errors.WithStack(err)
	return
}
