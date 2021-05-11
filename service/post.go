package service

import (
	"clover/model/mysql"
	"clover/model/redis"
	"clover/pkg/log"
	"clover/pkg/snowflake"
	"errors"
	"strconv"

	gormMysql "github.com/go-sql-driver/mysql"
)

var ErrorPostExisted = errors.New("post existed")
var ErrorPostUnexisted = errors.New("post unexisted")
var ErrorPostEmpty = errors.New("post list empty")

func CreatePost(post *mysql.PostCreateReq, author int64) (*mysql.Post, error) {

	postId := snowflake.GenSnowflakeID()
	p := &mysql.Post{
		PostID:      int64(postId),
		Title:       post.Title,
		Content:     post.Content,
		AuthorID:    author,
		CommunityID: post.CommunityID,
	}

	err := p.Insert()
	if sqlError, ok := err.(*gormMysql.MySQLError); ok {
		if sqlError.Number == 1062 {
			log.WithCategory("service.post").WithError(sqlError).Info("CreatePost: dumplicate post")
			return nil, ErrorPostExisted
		}
	}
	if err != nil {
		log.WithCategory("service.post").WithError(err).Error("CreatePost: insert post error")
		return nil, ErrorDBHandle
	}

	// query authorName & communityName
	user, err := mysql.QueryUserByUserID(author)
	if err != nil {
		log.WithCategory("service.post").WithError(err).Error("CreatePost: query user error")
		return nil, ErrorDBHandle
	}

	community, err := mysql.QueryCommunityByID(post.CommunityID)
	if err != nil {
		log.WithCategory("service.post").WithError(err).Error("CreatePost: query community error")
		return nil, ErrorDBHandle
	}

	postInfo := &redis.PostInfo{
		PostID:        strconv.FormatUint(postId, 10),
		AuthorID:      strconv.FormatInt(author, 10),
		CommunityID:   strconv.FormatInt(post.CommunityID, 10),
		AuthorName:    user.UserName,
		CommunityName: community.CommunityName,
		Title:         post.Title,
		Summary:       GetSummaryFromContent(post.Content, 128),
	}

	err = redis.CreatePost(postInfo)
	if err != nil {
		log.WithCategory("service.post").WithError(err).Error("CreatePost: redis error")
		return nil, ErrorRedisHandle
	}

	return p, nil
}

func GetPostList(page *mysql.PostListPage) ([]mysql.PostItem, error) {

	items, err := mysql.QueryPostList(page)
	if items == nil || len(items) == 0 {
		log.WithCategory("service.post").Info("GetPostList: empty post list")
		return nil, ErrorPostEmpty
	}

	if err != nil {
		log.WithCategory("service.post").WithError(err).Error("GetPostList: query post list error")
		return nil, ErrorDBHandle
	}

	postList := make([]mysql.PostItem, len(items))
	for index, item := range items {
		//query authorName
		var (
			authorName    string
			communityName string
			result        mysql.PostItem
		)
		author, err := mysql.QueryUserByUserID(item.AuthorID)
		if err == nil {
			authorName = author.UserName
		}

		community, err := mysql.QueryCommunityByID(item.CommunityID)
		if err == nil {
			communityName = community.CommunityName
		}

		result = mysql.PostItem{
			PostID:        item.PostID,
			AuthorName:    authorName,
			CommunityName: communityName,
			Title:         item.Title,
			Content:       item.Content,
			Status:        item.Status,
		}

		postList[index] = result
	}

	return postList, nil
}

type PostDetail struct {
	mysql.Post
	AuthorName    string `json:"author_name"`
	CommunityName string `json:"community_name"`
}

func GetPostDetail(id int64) (*PostDetail, error) {

	p, err := mysql.QueryPost(id)
	if p == nil {
		log.WithCategory("service.post").Info("GetPostDetail: post not existed")
		return nil, ErrorPostUnexisted
	}

	if err != nil {
		log.WithCategory("service.post").WithError(err).Error("GetPostDetail: query error")
		return nil, ErrorDBHandle
	}

	var (
		authorName    string
		communityName string
	)
	author, err := mysql.QueryUserByUserID(p.AuthorID)
	if err == nil {
		authorName = author.UserName
	}

	community, err := mysql.QueryCommunityByID(p.CommunityID)
	if err == nil {
		communityName = community.CommunityName
	}

	return &PostDetail{
		Post:          *p,
		AuthorName:    authorName,
		CommunityName: communityName,
	}, nil
}

func GetSummaryFromContent(content string, summaryLen int) string {
	if len(content) <= summaryLen {
		return content
	}

	return content[:summaryLen-1]
}
