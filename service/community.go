package service

import (
	"clover/model/mysql"
	"clover/pkg/log"
	"clover/pkg/snowflake"
	"errors"
)

var ErrorCommunityExisted = errors.New("community existed")
var ErrorCommunityEmpty = errors.New("community empty")

func CreateCommunity(c *mysql.CommunityCreateReq) (*mysql.Community, error) {

	//query exist
	community, err := mysql.QueryCommunity(c.CommunityName)
	if community != nil {
		log.WithCategory("service.community").Infof("CreateCommunity: community %s existed", c.CommunityName)
		return nil, ErrorCommunityExisted
	}

	if err != nil {
		log.WithCategory("service.community").WithError(err).Error("CreateCommunity: query error")
		return nil, ErrorDBHandle
	}

	// generate community id
	id := snowflake.GenSnowflakeID()
	community = &mysql.Community{
		CommunityID:   int64(id),
		CommunityName: c.CommunityName,
		Introduction:  c.Introduction,
	}

	err = community.Insert()
	if err != nil {
		log.WithCategory("service.community").WithError(err).Error("CreateCommunity: insert new community error")
		return nil, err
	}

	return community, nil

}

func GetCommunityList() (list []mysql.CommunityItem, err error) {

	list, err = mysql.QueryAllCommunities()
	if list == nil || len(list) == 0 {
		log.WithCategory("service.community").Info("GetCommunityList: empty community")
		err = ErrorCommunityEmpty
		return
	}

	if err != nil {
		log.WithCategory("service.community").WithError(err).Error("GetCommunityList: query error")
		err = ErrorDBHandle
	}

	return
}

func GetCommunityDetail(id int64) (community *mysql.Community, err error) {
	community, err = mysql.QueryCommunityByID(id)
	if community == nil {
		log.WithCategory("service.community").Info("GetCommunityDetail: empty community")
		err = ErrorCommunityEmpty
		return
	}

	if err != nil {
		log.WithCategory("service.community").WithError(err).Error("GetCommunityDetail: query error")
		err = ErrorDBHandle
	}

	return
}
