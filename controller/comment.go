package controller

import (
	"clover/model/mysql"
	"clover/pkg/log"
	"clover/service"
	"errors"

	"github.com/gin-gonic/gin"
)

func CommentCreateHandler(c *gin.Context) {

	var (
		req mysql.CommentCreateReq
		err error
	)

	err = c.ShouldBind(&req)
	if err != nil {
		log.WithCategory("controller.comment").WithError(err).Info("CommentCreateHandler: params invalid")
		Error400(c, CodeParamInvalid)
		return
	}

	userId, ok := c.Get(ContextKeyUserID)
	if !ok {
		log.WithCategory("controller.comment").WithError(err).Info("CommentCreateHandler: userid error")
		Error400(c, CodeTokenInvalid)
		return
	}

	err = service.CreateComment(userId.(int64), &req)
	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	Succ(c, nil)
}

func CommentListHandler(c *gin.Context) {

	postId, ok := c.Params.Get("post_id")
	if !ok {
		log.WithCategory("controller.comment").Info("CommentListHandler: invalid postId")
		Error400(c, CodeParamInvalid)
		return
	}

	list, err := service.GetCommentList(postId)
	if errors.Is(err, service.ErrorCommentEmpty) {
		Error(c, CodeCommentEmpty)
		return
	}

	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	Succ(c, list)
}
