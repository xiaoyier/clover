package controller

type RequestError struct {
	err  error
	code uint16
}

type ResponseCode uint16

const (
	CodeSucc ResponseCode = 1000 + iota
	CodeParamInvalid
	CodeUserExist
	CodeUserNotExist
	CodePasswordWrong
	CodeCommunityExist
	CodeCommunityEmpty
	CodeTokenInvalid
	CodeServerBusy
)

var CodeMessageMap = map[ResponseCode]string{
	CodeSucc:           "success",
	CodeParamInvalid:   "invalid params",
	CodeUserExist:      "existed user",
	CodeUserNotExist:   "user is not exist",
	CodePasswordWrong:  "password is wrong",
	CodeCommunityExist: "existed community",
	CodeCommunityEmpty: "empty community",
	CodeTokenInvalid:   "invalid token",
	CodeServerBusy:     "server busy",
}

func (c ResponseCode) Message() string {
	return CodeMessageMap[c]
}
