package controller

import (
	"clover/model/redis"
	"clover/pkg/log"
	"errors"

	"github.com/gin-gonic/gin"
)

type VoteReq struct {
	PostId    int64   `json:"post_id,string" binding:"required"`
	Direction float64 `json:"direction" binding:"required,oneOf=0 1 -1"`
}

func VoteHandler(c *gin.Context) {

	var (
		req VoteReq
		err error
	)
	err = c.ShouldBind(&req)
	if err != nil {
		log.WithCategory("controller.vote").WithError(err).Info("VoteHandler: params invalid")
		Error400(c, CodeParamInvalid)
		return
	}

	err = redis.Vote(req.PostId, req.Direction)
	if errors.Is(err, redis.ErrorVoteExpire) {
		Error(c, CodeVoteExpired)
	}

	if errors.Is(err, redis.ErrorVoted) {
		Error(c, CodeVoted)
		return
	}

	if err != nil {
		Error(c, CodeServerBusy)
		return
	}

	Succ(c, nil)
}
