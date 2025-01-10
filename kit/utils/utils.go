package utils

import "github.com/google/uuid"

func GenerateUUID() (string, error) {
	uuidObject, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuidObject.String(), nil
}
