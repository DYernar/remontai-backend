package util

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GetUsernameFromEmail(email string) string {
	if !strings.Contains(email, "@") {
		return email
	}
	return email[:strings.Index(email, "@")]
}

var appSecretKey = []byte("AppSecretKey")

func GenerateJWT(email string, id string) (string, error) {
	accessTokenExp := time.Now().Add(10000 * time.Second).Unix()
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["id"] = id
	accessTokenClaims["email"] = email
	accessTokenClaims["iat"] = time.Now().Unix()
	accessTokenClaims["exp"] = accessTokenExp
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	return accessToken.SignedString([]byte(appSecretKey))
}
