package usecase

type RegisterNewUserInput struct {
	PhoneNumber string
	FullName    string
	Password    string
}

type RegisterNewUserOutput struct {
	Id                  int64
	IsPhoneNumberExists bool
}

type LoginInput struct {
	PhoneNumber string
	Password    string
}

type LoginOutput struct {
	IsDataNotFound  bool
	IsPasswordWrong bool
	Token           string
}

type GetUserDataInput struct {
	Id int64
}

type GetUserDataOutput struct {
	PhoneNumber string
	FullName    string
}

type UpdateUserDataInput struct {
	Id          int64
	PhoneNumber string
	FullName    string
}

type UpdateUserDataOutput struct {
	Id                  int64
	IsPhoneNumberExists bool
}
