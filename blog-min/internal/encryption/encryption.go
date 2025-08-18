package encryption

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) string {
	saltedPassword := append(salt, []byte(password)...)
	hash := sha256.Sum256(saltedPassword)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func GenerateNewPassword(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}
	hashedPassword := hashPassword(password, salt)
	return hashedPassword, base64.StdEncoding.EncodeToString(salt), nil
}

func GeneratePasswordFromSalt(password, salt string) (string, error) {
	decodedSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	hashedPassword := hashPassword(password, decodedSalt)
	return hashedPassword, nil
}

func CheckPassword(loginPassword, dbPassword, salt string) bool {
	loginHash, _ := GeneratePasswordFromSalt(loginPassword, salt)
	return loginHash == dbPassword
}
