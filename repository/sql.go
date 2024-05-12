package repository

var (
	InsertNewUserQuery = `INSERT INTO users(phone_number, full_name, password) values ($1, $2, $3) returning id`

	UpdateUserDataQuery = `UPDATE users 
	set phone_number = $2, 
	full_name = $3,
	updated_at = now(),
	updated_by = $4
	WHERE id = $1`

	GetPasswordByPhoneNumberQuery = `SELECT id, password, phone_number FROM users WHERE phone_number = $1`

	UpdateTotalLoginById = `UPDATE users
	SET total_login = total_login + 1
	WHERE id = $1`

	GetUserDataByIdQuery = `SELECT id, full_name, phone_number FROM users WHERE id = $1`
)
