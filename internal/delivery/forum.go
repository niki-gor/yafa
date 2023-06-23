package delivery

import (
	"net/http"
	"strconv"
	"yafa/internal/model"
	"yafa/internal/repository"

	"github.com/labstack/echo/v4"
)

type ForumHandlersI interface {
	Create(ctx echo.Context) error
	Details(ctx echo.Context) error
	CreateThread(ctx echo.Context) error
	ForumThreads(ctx echo.Context) error
	ForumUsers(ctx echo.Context) error
}

type forumH struct {
	forumRepo  repository.ForumRepo
	userRepo   repository.UserRepo
	threadRepo repository.ThreadRepo
}

func NewForumHandler(f repository.ForumRepo, u repository.UserRepo, t repository.ThreadRepo) ForumHandlersI {
	return &forumH{
		forumRepo:  f,
		userRepo:   u,
		threadRepo: t,
	}
}

func (h *forumH) Create(ctx echo.Context) error {
	forum := model.Forum{}
	err := ctx.Bind(&forum)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	checkForum, err := h.forumRepo.GetBySlug(forum.Slug)
	if err == nil {
		return ctx.JSON(http.StatusConflict, checkForum)
	}

	checkUser, err := h.userRepo.GetByNickname(forum.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	forum.User = checkUser.Nickname

	newForum, err := h.forumRepo.Create(forum)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, newForum)
}

func (h *forumH) Details(ctx echo.Context) error {
	checkForum, err := h.forumRepo.GetBySlug(ctx.Param("slug"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, checkForum)
}

func (h *forumH) CreateThread(ctx echo.Context) error {
	checkForum, err := h.forumRepo.GetBySlug(ctx.Param("slug"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	thread := model.Thread{}
	err = ctx.Bind(&thread)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if thread.Slug != "" {
		checkThread, err := h.threadRepo.GetBySlugOrId(thread.Slug)
		if err == nil {
			return ctx.JSON(http.StatusConflict, checkThread)
		}
	}

	thread.Forum = checkForum.Slug

	checkAuthor, err := h.userRepo.GetByNickname(thread.Author)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	thread.Author = checkAuthor.Nickname

	newThread, err := h.threadRepo.Create(thread)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, newThread)
}

func (h *forumH) ForumThreads(ctx echo.Context) error {
	slug := ctx.Param("slug")
	_, err := h.forumRepo.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
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

	threads, err := h.forumRepo.GetThreads(slug, since, limit, desc)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, threads)
}

func (h *forumH) ForumUsers(ctx echo.Context) error {
	slug := ctx.Param("slug")
	forum, err := h.forumRepo.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
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

	users, err := h.forumRepo.GetUsers(forum, since, limit, desc)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return ctx.JSON(http.StatusOK, users)
}
