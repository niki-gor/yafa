package router

import (
	"yafa/internal/delivery"
	"yafa/internal/pkg/serializer"
	"yafa/internal/repository"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
)

func New(db *pgx.ConnPool) *echo.Echo {
	e := echo.New()

	userRepo := repository.NewUserRepo(db)
	forumRepo := repository.NewForumRepo(db)
	threadRepo := repository.NewThreadRepo(db)
	postRepo := repository.NewPostRepo(db)
	serviceRepo := repository.NewServiceRepo(db)

	userHandler := delivery.NewUserHandler(userRepo)
	forumHandler := delivery.NewForumHandler(forumRepo, userRepo, threadRepo)
	threadHandler := delivery.NewThreadHandler(threadRepo, userRepo)
	postHandler := delivery.NewPostHandler(postRepo)
	serviceHandler := delivery.NewServiceHandler(serviceRepo)

	e.POST("/api/forum/create", forumHandler.Create)
	e.GET("/api/forum/:slug/details", forumHandler.Details)
	e.POST("/api/forum/:slug/create", forumHandler.CreateThread)
	e.GET("/api/forum/:slug/threads", forumHandler.ForumThreads)
	e.GET("/api/forum/:slug/users", forumHandler.ForumUsers)

	e.GET("/api/post/:id/details", postHandler.GetDetails)
	e.POST("/api/post/:id/details", postHandler.UpdateDetails)

	e.GET("/api/service/status", serviceHandler.Status)
	e.POST("/api/service/clear", serviceHandler.Clear)

	e.POST("/api/thread/:slug_or_id/create", threadHandler.CreatePost)
	e.POST("/api/thread/:slug_or_id/vote", threadHandler.CreateVote)
	e.GET("/api/thread/:slug_or_id/details", threadHandler.Details)
	e.GET("/api/thread/:slug_or_id/posts", threadHandler.ThreadPost)
	e.POST("/api/thread/:slug_or_id/details", threadHandler.Update)

	e.POST("/api/user/:nickname/create", userHandler.Create)
	e.GET("/api/user/:nickname/profile", userHandler.GetByNickname)
	e.POST("/api/user/:nickname/profile", userHandler.Update)

	e.JSONSerializer = serializer.EasyJsonSerializer{}

	return e
}
