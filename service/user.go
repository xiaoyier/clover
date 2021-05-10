package service

import (
	"clover/model"
	"clover/pkg/jwt"
	"clover/pkg/log"
	"clover/pkg/snowflake"
	md5 "crypto/md5"
	"encoding/hex"
	"errors"
)

var ErrorDBHandle = errors.New("mysql handle error")
var ErrorUserExist = errors.New("user exists")
var ErrorUserNotExist = errors.New("user is not existed")
var ErrorPasswordWrong = errors.New("wrong password")

func CreateUser(u *model.UserRegisterReq) error {

	user, err := model.QueryUserByUserName(u.UserName)
	if err != nil {
		log.WithCategory("service.user").WithError(err).Error("CreateUser: query user failed")
		return ErrorDBHandle
	}

	if user != nil {
		log.WithCategory("service.user").Info("CreateUser: user existed")
		return ErrorUserExist
	}

	// Gen User ID
	userId := snowflake.GenSnowflakeID()
	if userId == 0 {
		log.WithCategory("service.user").Error("CreateUser: generate user id failed")
		return ErrorDBHandle
	}

	user = &model.User{
		UserName:   u.UserName,
		UserPasswd: MD5Encrypt(u.Password),
		UserID:     int64(userId),
	}

	err = user.Insert()
	if err != nil {
		log.WithCategory("service.user").WithError(err).Error("CreateUser: insert user failed")
		return ErrorDBHandle
	}
	return nil
}

func Login(u *model.UserLoginReq) (rsp *model.UserLoginRsp, err error) {

	// query if not exist user
	user, err := model.QueryUserByUserName(u.UserName)
	if user == nil {
		log.WithCategory("service.user").Info("Login: user not existed: ", u.UserName)
		return rsp, ErrorUserNotExist
	}

	if err != nil {
		log.WithCategory("service.user").WithError(err).Error("Login: query user failed")
		return rsp, ErrorDBHandle
	}

	// check password
	if user.UserPasswd != MD5Encrypt(u.Password) {
		log.WithCategory("service.user").Info("Login: password wrong.")
		return rsp, ErrorPasswordWrong
	}

	// generate jwt token
	login, refresh, err := jwt.GenToken(user.UserID)
	if err != nil {
		log.WithCategory("service.user").WithError(err).Error("Login: generate token failed")
		return rsp, err
	}

	return &model.UserLoginRsp{
		UserID:       user.UserID,
		UserName:     user.UserName,
		LoginToken:   login,
		RefreshToken: refresh,
	}, nil
}

func MD5Encrypt(src string) string {

	sha := md5.New()
	sha.Write([]byte(src))
	return hex.EncodeToString(sha.Sum(nil))
}
