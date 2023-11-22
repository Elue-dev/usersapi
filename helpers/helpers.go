package helpers

import (
	"github.com/elue-dev/usersapi/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func ComparePasswordWithHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func DatabaseUserToUserModel (dbUser models.User) models.CustomUser {
	return models.CustomUser{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		DeletedAt:    dbUser.DeletedAt,
		FirstName:    dbUser.FirstName,
		LastName:     dbUser.LastName,
		Email:        dbUser.Email,
		Password:     dbUser.Password,
		Avatar:       dbUser.Avatar,
	}
}