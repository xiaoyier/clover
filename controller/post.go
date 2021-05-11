package controller

import (
	"clover/model/mysql"
	"clover/pkg/log"
	"clover/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostCreateHandler(c *gin.Context) {
	var (
		req mysql.PostCreateReq
		rsp *mysql.Post
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		log.WithCategory("post").WithError(err).Info("PostCreateHandler: params invalid")
		Error400(c, CodeParamInvalid)
		return
	}

	author, exists := c.Get(ContextKeyUserID)
	if !exists {
		log.WithCategory("post").Info("PostCreateHandler: not authorized")
		Error401(c, CodeTokenInvalid)
		return
	}

	authorId, ok := author.(int64)
	if !ok {
		log.WithCategory("post").Info("PostCreateHandler: not authorized")
		Error401(c, CodeTokenInvalid)
		return
	}

	rsp, err = service.CreatePost(&req, authorId)
	if errors.Is(err, service.ErrorPostExisted) {
		Error(c, CodePostExisted)
		return
	}

	if err != nil {
		Error(c, CodeServerBusy)
		return
	}

	Succ(c, rsp)
}

func PostDetailHandler(c *gin.Context) {

	postId, ok := c.Params.Get("id")
	if !ok {
		log.WithCategory("controller.post").Info("PostDetailHandler: post_id error")
		Error400(c, CodeParamInvalid)
		return
	}

	id, err := strconv.ParseInt(postId, 10, 64)
	if id == 0 || err != nil {
		log.WithCategory("controller.post").WithError(err).Error("PostDetailHandler: parse post id error")
		Error400(c, CodeParamInvalid)
		return
	}

	rsp, err := service.GetPostDetail(id)
	if errors.Is(err, service.ErrorPostUnexisted) {
		Error(c, CodePostUnexist)
		return
	}

	if err != nil {
		Error(c, CodeServerBusy)
		return
	}

	Succ(c, rsp)
}

func PostListHandler(c *gin.Context) {

	var (
		page mysql.PostListPage
		err  error
		rsp  []mysql.PostItem
	)

	page = mysql.PostListPage{
		PageSize:   20,
		PageNumber: 0,
	}
	err = c.BindQuery(&page)
	if err != nil {
		log.WithCategory("post").WithError(err).Info("PostListHandler: invalid params")
		Error400(c, CodeParamInvalid)
		return
	}

	rsp, err = service.GetPostList(&page)
	if errors.Is(err, service.ErrorPostEmpty) {
		Error(c, CodePostEmpty)
		return
	}

	if err != nil {
		Error(c, CodeServerBusy)
		return
	}

	Succ(c, rsp)
}

func PostListV2Handler(c *gin.Context) {

}
