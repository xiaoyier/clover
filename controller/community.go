package controller

import (
	"clover/model"
	"clover/pkg/log"
	"clover/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CommunityCreateHandler(c *gin.Context) {

	var (
		err error
		req model.CommunityCreateReq
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.WithCategory("controller.community").WithError(err).Info("CreateCommunityHandler: params error")
		Error400(c, CodeParamInvalid)
		return
	}

	rsp, err := service.CreateCommunity(&req)
	if errors.Is(err, service.ErrorCommunityExisted) {
		Error(c, CodeCommunityExist)
		return
	}

	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	data := struct {
		CommunityID   int64  `json:"community_id,string"`
		CommunityName string `json:"community_name"`
		Introduction  string `json:"introduction"`
	}{
		CommunityID:   rsp.CommunityID,
		CommunityName: rsp.CommunityName,
		Introduction:  rsp.Introduction,
	}

	Succ(c, data)
}

func CommunityListHandler(c *gin.Context) {

	rsp, err := service.GetCommunityList()
	if errors.Is(err, service.ErrorCommunityEmpty) {
		Error(c, CodeCommunityEmpty)
		return
	}

	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	Succ(c, rsp)
}

func CommunityDetailHandler(c *gin.Context) {

	id, ok := c.Params.Get("id")
	if !ok {
		log.WithCategory("controller.community").Info("CommunityDetailHandler: error community id")
		Error400(c, CodeParamInvalid)
		return
	}

	communityId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithCategory("controller.community").WithError(err).Error("CommunityDetailHandler: parse community id error")
		Error(c, CodeParamInvalid)
		return
	}

	rsp, err := service.GetCommunityDetail(communityId)
	if errors.Is(err, service.ErrorCommunityEmpty) {
		Error(c, CodeCommunityEmpty)
		return
	}

	if err != nil {
		Error500(c, CodeServerBusy)
		return
	}

	Succ(c, rsp)
}
