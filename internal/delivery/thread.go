package delivery

import (
	"net/http"
	"strconv"
	"yafa/internal/model"
	"yafa/internal/repository"

	"github.com/labstack/echo/v4"
)

type ThreadHandler interface {
	CreatePost(ctx echo.Context) error
	CreateVote(ctx echo.Context) error
	Details(ctx echo.Context) error
	ThreadPost(ctx echo.Context) error
	Update(ctx echo.Context) error
}

type threadHandler struct {
	threadRepo repository.ThreadRepo
	userRepo   repository.UserRepo
}

func NewThreadHandler(t repository.ThreadRepo, u repository.UserRepo) ThreadHandler {
	return &threadHandler{threadRepo: t, userRepo: u}
}

func (h *threadHandler) CreatePost(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")

	thread, err := h.threadRepo.GetBySlugOrId(slugOrId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var Posts []model.Post
	err = ctx.Bind(&Posts)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	posts := model.Posts{Posts: Posts}
	if len(posts.Posts) == 0 {
		return ctx.JSON(http.StatusCreated, posts.Posts)
	}

	for _, item := range posts.Posts {
		_, err := h.userRepo.GetByNickname(item.Author)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		if item.Parent != 0 {
			err = h.threadRepo.CheckPost(item.Parent, thread.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusConflict)
			}
		}
	}

	response, err := h.threadRepo.CreatePosts(thread, posts)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict)
	}

	return ctx.JSON(http.StatusCreated, response.Posts)
}

func (h *threadHandler) CreateVote(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")

	thread, err := h.threadRepo.GetBySlugOrId(slugOrId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var vote model.VoteRequest
	err = ctx.Bind(&vote)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	checkUser, err := h.userRepo.GetByNickname(vote.Nickname)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	vote1, err := h.threadRepo.CheckVotes(checkUser.Id, thread.Id)
	if err == nil && vote.Voice == vote1.Voice {
		return ctx.JSON(http.StatusOK, thread)
	}
	if err != nil {
		err = h.threadRepo.CreateVote(checkUser.Id, vote, thread)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		thread.Votes += vote.Voice
	} else {
		_, err = h.threadRepo.UpdateVote(vote, vote1.Id)
		if err == nil {
			thread.Votes += 2 * vote.Voice
		}
	}
	return ctx.JSON(http.StatusOK, thread)
}

func (h *threadHandler) Details(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")

	thread, err := h.threadRepo.GetBySlugOrId(slugOrId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (h *threadHandler) ThreadPost(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")

	thread, err := h.threadRepo.GetBySlugOrId(slugOrId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	sort := ctx.FormValue("sort")
	if sort == "" {
		sort = "flat"
	}

	limit := 0
	if limitVal := ctx.FormValue("limit"); limitVal != "" {
		limit, err = strconv.Atoi(limitVal)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	} else {
		limit = 100
	}

	since := ctx.FormValue("since")

	desc, _ := strconv.ParseBool(ctx.FormValue("desc"))

	posts, err := h.threadRepo.GetThreadPosts(thread, since, sort, limit, desc)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, posts)
}

func (h *threadHandler) Update(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")

	thread, err := h.threadRepo.GetBySlugOrId(slugOrId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var updateThread model.Thread
	err = ctx.Bind(&updateThread)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if updateThread.Title == "" && updateThread.Message == "" {
		return ctx.JSON(http.StatusOK, thread)
	}
	if updateThread.Title == "" {
		updateThread.Title = thread.Title
	}
	if updateThread.Message == "" {
		updateThread.Message = thread.Message
	}

	thread, err = h.threadRepo.Update(thread, updateThread)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, thread)
}
