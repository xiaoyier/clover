package service

import (
	"clover/model/mysql"
	"clover/pkg/log"
	"clover/pkg/snowflake"
	"errors"
	"strconv"
)

var ErrorCommentEmpty = errors.New("empty comment")

func CreateComment(authorId int64, comment *mysql.CommentCreateReq) error {

	commentId := snowflake.GenSnowflakeID()
	model := &mysql.Comment{
		CommentID: int64(commentId),
		PostID:    comment.PostID,
		AuthorID:  authorId,
		Content:   comment.Content,
	}

	return model.Insert()
}

func GetCommentList(postId string) ([]mysql.Comment, error) {

	id, _ := strconv.ParseInt(postId, 10, 64)
	list, err := mysql.QueryCommentList(id)
	if list == nil || len(list) == 0 {
		log.WithCategory("service.comment").Info("GetCommentList: empty comment")
		return nil, ErrorCommentEmpty
	}

	if err != nil {
		log.WithCategory("service.comment").WithError(err).Error("GetCommentList: query error")
		return nil, err
	}

	return list, nil
}
