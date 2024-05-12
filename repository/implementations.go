package repository

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *Repository) InsertNewUser(ctx context.Context, input InsertNewUserInput) (InsertNewUserOutput, error) {
	var (
		output InsertNewUserOutput
		newId  int64
	)

	err := r.Db.QueryRowContext(ctx, InsertNewUserQuery, input.PhoneNumber, input.FullName, input.Password).Scan(&newId)
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
	_, err := r.Db.ExecContext(ctx, UpdateUserDataQuery, input.Id, input.PhoneNumber, input.FullName, input.Id)
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
	err = r.Db.QueryRowContext(ctx, GetPasswordByPhoneNumberQuery, input.PhoneNumber).Scan(&output.Id, &output.Password, &output.PhoneNumber)
	err = errors.WithStack(err)
	return
}

func (r *Repository) UpdateTotalLoginById(ctx context.Context, input UpdateTotalLoginByIdInput) (err error) {
	_, err = r.Db.ExecContext(ctx, UpdateTotalLoginById, input.Id)

	err = errors.WithStack(err)
	return err
}

func (r *Repository) GetUserDataById(ctx context.Context, input GetUserDataByIdInput) (output GetUserDataByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, GetUserDataByIdQuery, input.Id).Scan(&output.Id, &output.FullName, &output.PhoneNumber)
	err = errors.WithStack(err)
	return
}
