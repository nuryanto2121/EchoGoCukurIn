package util

import "golang.org/x/crypto/bcrypt"

// Hash :
func Hash(text string) (string, error) {
	pwd := []byte(text)

	hashedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return text, err
	}
	return string(hashedPwd), nil
}

// Compare :
func Compare(hashed string, text string) (bool, error) {
	hashedPwd := []byte(hashed)
	pwd := []byte(text)

	err := bcrypt.CompareHashAndPassword(hashedPwd, pwd)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetPassword(password string) []byte {
	return []byte(password)
}

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}
