package utils

import (
	"os"
	"time"
	"uas/app/models"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(user models.User) (string, error) {
	claims := models.JWTClaims{
		UserID: user.ID,
		Username: user.Username,
		RoleName: user.RoleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func GenerateRefreshToken(user models.User) (string, error) {
    claims := jwt.MapClaims{
        "userId": user.ID,
        "role":   user.RoleName,
        "exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
        "type":   "refresh",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JwtSecret)
}

func ValidateToken(tokenString string) (*models.JWTClaims, error) { 
    token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{},func(token *jwt.Token) (interface {}, error) { 
        return JwtSecret, nil 
    }) 
 
    if err != nil { 
        return nil, err 
    } 
 
    if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid { 
        return claims, nil 
    } 
 
    return nil, jwt.ErrInvalidKey 
}