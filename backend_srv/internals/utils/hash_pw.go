package utils

import ("golang.org/x/crypto/bcrypt")

func HashPassword(password string) (string, error) {
	// use bcrypt
	byteResult, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(byteResult), nil

}

type ErrWrongPasswordType struct{}

func (e ErrWrongPasswordType) Error() string {
	return "password tidak sesuai"
}

var ErrWrongPassword ErrWrongPasswordType

func ComparePassword(hashedPassword string, inputPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}
