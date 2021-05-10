package controller

import (
	"clover/model"
	"clover/pkg/jwt"
	"clover/pkg/log"
	"clover/service"
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserSignUpHandler(c *gin.Context) {
	var (
		err error
		req model.UserRegisterReq
	)

	err = c.ShouldBind(&req)
	if err != nil {
		log.WithCategory("controller.user").WithError(err).Info("UserRegisterHandler: params error")
		Error400(c, CodeParamInvalid)
		return
	}
	//
	//if !CheckInvalidChar(req.UserName) || !CheckInvalidChar(req.Password) {
	//	Error400(c, CodeParamInvalid)
	//	return
	//}

	err = service.CreateUser(&req)
	if errors.Is(err, service.ErrorUserExist) {
		Error(c, CodeUserExist)
		return
	}

	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	Succ(c, nil)
}

func UserLoginHandler(c *gin.Context) {
	var (
		err error
		req model.UserLoginReq
	)

	err = c.ShouldBind(&req)
	if err != nil {
		log.WithTraceID("controller.user").WithError(err).Info("UserLoginHandler: params error")
		Error400(c, CodeParamInvalid)
		return
	}

	rsp, err := service.Login(&req)
	if errors.Is(err, service.ErrorUserNotExist) {
		Error(c, CodeUserNotExist)
		return
	}
	if errors.Is(err, service.ErrorPasswordWrong) {
		Error(c, CodePasswordWrong)
		return
	}

	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	Succ(c, rsp)
}

func UserRefreshTokenHandler(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if !jwt.CheckValid(refreshToken) {
		log.WithCategory("controller.user").Info("UserRefreshTokenHandler: invalid refresh token")
		Error400(c, CodeTokenInvalid)
		return
	}

	accessToken := c.Request.Header.Get("Authorization")
	contents := strings.Split(accessToken, " ")
	if !(len(contents) == 2 && contents[0] == "Bearer") {
		log.WithCategory("controller.user").Info("UserRefreshTokenHandler: Token format error")
		Error400(c, CodeTokenInvalid)
		return
	}

	userId, err := jwt.ParseUserID(contents[1])
	if err != nil || userId == 0 {
		log.WithCategory("controller.user").WithError(err).Error("UserRefreshTokenHandler: parse userid from token failed")
		Error400(c, CodeTokenInvalid)
		return
	}

	loginToken, refreshToken, err := jwt.GenToken(userId)
	if err != nil {
		log.WithCategory("controller.user").WithError(err).Error("UserRefreshTokenHandler: gen token failed")
		Error400(c, CodeServerBusy)
		return
	}

	data := struct {
		RefreshToken string `json:"refresh_token"`
		LoginToken   string `json:"login_token"`
	}{RefreshToken: refreshToken, LoginToken: loginToken}

	Succ(c, data)
}

func CheckInvalidChar(param string) bool {

	reg, err := regexp.Compile(`[\\ud83c\\udc00-\\ud83c\\udfff]|[\\ud83d\\udc00-\\ud83d\\udfff]|[\\u2600-\\u27ff]`)
	if err != nil {
		log.WithCategory("controller.user").WithError(err).Error("checkInvalidChar: get reg error")
		return true
	}

	return reg.Match([]byte(param))
}
