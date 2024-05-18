package controllers

import (
	"golang-template/dto"
	"golang-template/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userController struct {
	UserService  services.UserService
	TokenService services.TokenService
}

func NewUserController() *userController {
	return &userController{
		UserService: services.NewUserService(),
	}
}

func (u *userController) GetUserByToken(c echo.Context) error {
	user, err := u.TokenService.UserByToken(c)
	if err != nil {
		return echo.ErrUnauthorized
	} else {
		return c.JSON(http.StatusOK, user)
	}
}

func (u *userController) UpdateUser(c echo.Context) error {
	updateForm := &dto.RegisterForm{}
	err := c.Bind(updateForm)
	if err != nil {
		return c.String(http.StatusBadRequest, "All user fields must be provided!")
	}

	updatedUser, err := u.UserService.UpdateUser(updateForm, c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedUser)
}

func (u *userController) DeleteUser(c echo.Context) error {
	err := u.UserService.DeleteUser(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in removing user")
	}

	return c.String(http.StatusOK, "account removed successfully")
}

func (u *userController) Logout(c echo.Context) error {
	err := u.UserService.Logout(c)
	if err != nil {
		return err
	} else {
		return c.String(http.StatusOK, "Logout was successful")
	}
}
