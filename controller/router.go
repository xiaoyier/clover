package controller

import (
	"clover/model/mysql"
	"clover/pkg/jwt"
	"clover/pkg/limiter"
	"clover/pkg/log"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-uuid"

	"github.com/gin-gonic/gin"
)

var accessLog = flag.String("access_log_file", "", "the file record the gin access log(full path)")
var errorLog = flag.String("error_log_file", "", "the file record the gin error log(full path)")

func InitRouter() *gin.Engine {

	accessWriter := os.Stdout
	errWriter := os.Stderr
	if len(*accessLog) > 0 {
		if file, err := os.Open(*accessLog); err == nil {
			accessWriter = file
		}
	}

	if len(*errorLog) > 0 {
		if file, err := os.Open(*errorLog); err == nil {
			errWriter = file
		}
	}

	router := gin.New()
	router.Use(Limit(1))
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: accessWriter,
	}))
	router.Use(gin.RecoveryWithWriter(errWriter))
	router.Use(Trace)

	router.GET("/hello", HelloHandler)

	v1Group := router.Group("/api/v1")
	v1Group.POST("/user/signup", UserSignUpHandler)
	v1Group.POST("/user/login", UserLoginHandler)
	v1Group.GET("/user/token", UserRefreshTokenHandler)

	v1Group.Use(JWT)
	{
		v1Group.POST("/community/create", CommunityCreateHandler)
		v1Group.GET("/community/list", CommunityListHandler)
		v1Group.GET("/community/:id", CommunityDetailHandler)

		v1Group.POST("/post", PostCreateHandler)
		v1Group.GET("/post/:id", PostDetailHandler)
		v1Group.GET("/post/list", PostListHandler)
		v1Group.GET("/post/list2", PostListV2Handler)
		v1Group.GET("/post/list2/:community_id", PostCommunityListHandler)

		v1Group.POST("/vote", VoteHandler)

		v1Group.POST("/comment", CommentCreateHandler)
		v1Group.GET("/comment/:post_id", CommentListHandler)
	}

	router.NoRoute(func(context *gin.Context) {
		context.String(http.StatusNotFound, "hello boy")
	})

	return router
}

const ContextKeyUserID = "_user_id"

func JWT(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if len(token) == 0 {
		log.WithCategory("middleware.jwt").Info("JWT: authorization token empty")
		c.Abort()
		Error401(c, CodeTokenInvalid)
		return
	}
	contents := strings.Split(token, " ")
	if !(len(contents) == 2 && contents[0] == "Bearer") {
		log.WithCategory("middleware.jwt").Info("JWT: Token format error")
		c.Abort()
		Error401(c, CodeTokenInvalid)
		return
	}

	userId, err := jwt.ParseUserID(contents[1])
	if err != nil {
		log.WithCategory("middlwware.jwt").WithError(err).Error("parse userid from token failed")
		c.Abort()
		Error401(c, CodeTokenInvalid)
		return
	}

	if userId == 0 {
		log.WithCategory("middlwware.jwt").Info("JWT: invalid user id")
		c.Abort()
		Error401(c, CodeTokenInvalid)
		return
	}

	user, err := mysql.QueryUserByUserID(userId)
	if user == nil || err != nil {
		log.WithCategory("middleware.jwt").WithError(err).Error("JWT: user not login")
		c.Abort()
		Error401(c, CodeTokenInvalid)
		return
	}

	c.Set(ContextKeyUserID, userId)
	c.Next()
}

func Limit(capacity int64) gin.HandlerFunc {

	bucket := limiter.GetBucket()
	if bucket == nil {
		bucket = limiter.Init(capacity)
	}

	return func(context *gin.Context) {
		count := bucket.TakeAvailable(1)
		if count < 1 {
			log.WithCategory("controller.router").Error("Limit: limited")
			context.Abort()
			Error400(context, CodeRequestLimit)
			return
		}

		context.Next()
	}
}

const TraceIDKey = "_clover_trace"

func Trace(c *gin.Context) {

	traceId, _ := uuid.GenerateUUID()
	c.Set(TraceIDKey, traceId)
	c.Next()
}
