package service

type authenticator interface {
	CompareHashAndPassword(hash, password string) error
	GeneratePasswordHash(password string) (string, error)
}
