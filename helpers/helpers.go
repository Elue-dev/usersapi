package helpers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/elue-dev/usersapi/database"
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

func DatabaseUsersArrToUserModel (dbUsers []models.User) []models.CustomUser {
	users := []models.CustomUser{}

	for _, dbUser := range dbUsers {
		users = append(users, DatabaseUserToUserModel(dbUser))
	}
	return users
}

func GenerateToken(userID string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	expirationTime := time.Now().Add(24 * time.Hour) // 1 day

	claims := jwt.MapClaims{
		"user": userID,
		"exp":  expirationTime.Unix(),
	}

	// Create the JWT token with the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenFromHeaders(r *http.Request) string {
	headers := r.Header.Get("Authorization")
	tokenStr := strings.Split(headers, " ")
	return tokenStr[1]
}

func GetUserFromToken(r *http.Request) (models.User, error) {
	tokenString := getTokenFromHeaders(r)
	
    // Parse the JWT token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return models.User{}, err
    }

    // Check if the token is valid
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // Extract user information from the token
        userID, ok := claims["user"].(string)
        if !ok {
            return models.User{}, errors.New("invalid token format")
        }

        // Query the database for the user based on the extracted email
        var user models.User
        result := database.DB.Where("id = ?", userID).First(&user)

        if result.Error != nil {
            return models.User{}, result.Error
        }

        return user, nil
    }

    return models.User{}, errors.New("invalid token")
}


