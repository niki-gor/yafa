package delivery

import (
	"net/http"
	"yafa/internal/model"
	"yafa/internal/repository"

	"github.com/labstack/echo/v4"
)

type UserHandlersI interface {
	Create(ctx echo.Context) error
	GetByNickname(ctx echo.Context) error
	Update(ctx echo.Context) error
}

type userH struct {
	userRepo repository.UserRepo
}

func NewUserHandler(u repository.UserRepo) UserHandlersI {
	return &userH{
		userRepo: u,
	}
}

func (h *userH) Create(ctx echo.Context) error {
	var user model.User
	err := ctx.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	user.Nickname = ctx.Param("nickname")

	users, err := h.userRepo.GetByEmailOrNick(user.Email, user.Nickname)
	if err == nil && len(users) > 0 {
		return ctx.JSON(http.StatusConflict, users)
	}

	_, err = h.userRepo.Create(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (h *userH) GetByNickname(ctx echo.Context) error {
	user, err := h.userRepo.GetByNickname(ctx.Param("nickname"))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, user)
}

func (h *userH) Update(ctx echo.Context) error {
	newUserData, err := h.userRepo.GetByNickname(ctx.Param("nickname"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	err = ctx.Bind(&newUserData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	checkUser, _ := h.userRepo.GetByEmail(newUserData.Email)
	if checkUser.Nickname != "" && checkUser.Nickname != newUserData.Nickname {
		return echo.NewHTTPError(http.StatusConflict)
	}

	user, err := h.userRepo.Update(newUserData)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if newUserData.About == "" {
		user.About = newUserData.About
	}

	return ctx.JSON(http.StatusOK, user)
}
