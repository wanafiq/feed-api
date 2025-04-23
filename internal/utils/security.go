package utils

import "golang.org/x/crypto/bcrypt"

func Hash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyHash(hashedValue, value string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	return err == nil
}
