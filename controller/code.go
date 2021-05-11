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
	CodePostExisted
	CodePostUnexist
	CodePostEmpty
	CodeVoteExpired
	CodeVoted
	CodeCommentEmpty

	CodeTokenInvalid
	CodeRequestLimit
	CodeServerBusy
)

var CodeMessageMap = map[ResponseCode]string{
	CodeSucc:           "success",
	CodeParamInvalid:   "invalid params",
	CodeUserExist:      "existed user",
	CodeUserNotExist:   "user is not exist",
	CodePasswordWrong:  "password is wrong",
	CodeCommunityExist: "existed community",
	CodeCommunityEmpty: "empty community list",
	CodePostExisted:    "existed post",
	CodePostUnexist:    "unexist post",
	CodePostEmpty:      "empty post list",
	CodeVoteExpired:    "vote expired",
	CodeVoted:          "already voted",
	CodeCommentEmpty:   "empty comment",
	CodeTokenInvalid:   "invalid token",
	CodeRequestLimit:   "request limited",
	CodeServerBusy:     "server busy",
}

func (c ResponseCode) Message() string {
	return CodeMessageMap[c]
}
