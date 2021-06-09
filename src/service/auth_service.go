package service

type AuthService interface {
	Login(email string, password string) (string, error)
}
