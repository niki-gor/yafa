package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"yafa/internal/model"
	"yafa/internal/repository"

	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	GetDetails(ctx echo.Context) error
	UpdateDetails(ctx echo.Context) error
}

type postHandler struct {
	postRepo repository.PostRepo
}

func NewPostHandler(p repository.PostRepo) PostHandler {
	return &postHandler{
		postRepo: p,
	}
}

func (h *postHandler) GetDetails(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	related := strings.Split(ctx.FormValue("related"), ",")
	post, err := h.postRepo.Get(id, related)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, post)
}

func (h *postHandler) UpdateDetails(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var newPost model.Post
	err = ctx.Bind(&newPost)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var related []string
	postInfo, err := h.postRepo.Get(id, related)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	oldPost := postInfo.Post
	if newPost.Message == "" || oldPost.Message == newPost.Message {
		return ctx.JSON(http.StatusOK, oldPost)
	}

	post, err := h.postRepo.Update(id, newPost)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, post)
}
