package middleware

import (
	"golang-template/repository"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mid struct {
	userRepo repository.UserRepository
}

func NewMiddleware() *mid {
	userRepo := repository.NewUserRepositoryGORM()
	return &mid{
		userRepo: userRepo,
	}
}

func (m *mid) AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := GetToken(c)
		if token == "" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		uid, err := NewTokenRepositoryGORM().GetUserIdByToken(token)
		if err != nil {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		user, err := m.userRepo.GetUserById(uuid.MustParse(uid))
		if err != nil {
			response := map[string]interface{}{
				"code":    http.StatusNotFound,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusNotFound, response)
		}

		if user.User_role != "3" && user.User_role != "1" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		return next(c)
	}
}

func (m *mid) AuthenticateMitra(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := GetToken(c)
		if token == "" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		uid, err := NewTokenRepositoryGORM().GetUserIdByToken(token)
		if err != nil {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		user, err := m.userRepo.GetUserById(uuid.MustParse(uid))
		if err != nil {
			response := map[string]interface{}{
				"code":    http.StatusNotFound,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusNotFound, response)
		}

		if user.User_role != "2" && user.User_role != "1" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		return next(c)
	}
}

func (m *mid) AuthenticateAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := GetToken(c)
		if token == "" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		uid, err := NewTokenRepositoryGORM().GetUserIdByToken(token)
		if err != nil {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		user, err := m.userRepo.GetUserById(uuid.MustParse(uid))
		if err != nil {
			response := map[string]interface{}{
				"code":    http.StatusNotFound,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusNotFound, response)
		}

		if user.User_role != "1" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		return next(c)
	}
}
