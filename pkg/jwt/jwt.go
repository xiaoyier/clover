package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CloverClaims struct {
	jwt.StandardClaims
	UserID int64
}

const ExpireTimeLoginToken = time.Hour * 24 * 7
const ExpireTimeRefreshToken = time.Hour

var cloverSecret = []byte("golang 大法好!")

func GenToken(userID int64) (loginToken, refreshToken string, err error) {

	loginToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, CloverClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ExpireTimeLoginToken).Unix(),
			Issuer:    "clover",
		},
		userID,
	}).SignedString(cloverSecret)

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ExpireTimeRefreshToken).Unix(),
		Issuer:    "clover",
	}).SignedString(cloverSecret)

	return
}

func ParseUserID(loginToken string) (userId int64, err error) {

	var token *jwt.Token
	claims := &CloverClaims{}
	token, err = jwt.ParseWithClaims(loginToken, claims, func(token *jwt.Token) (interface{}, error) {
		return cloverSecret, nil
	})

	userId = claims.UserID
	if err != nil {
		return
	}

	if !token.Valid {
		err = errors.New("invalid token")
	}

	return
}

func CheckValid(refreshToken string) bool {
	token, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return cloverSecret, nil
	})

	return token != nil && token.Valid
}
