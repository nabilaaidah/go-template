package middleware

import (
	"errors"
	"golang-template/config"
	"golang-template/models"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TokenRepositoryGORM interface {
	GetUserIdByToken(token string) (string, error)
	SaveToken(user *models.User, token string) error
	SaveTokenResetEmail(user *models.User, token string) error
	GetTokenByToken(token string) (*models.Token, error)
	GetUserByToken(token string) (*models.User, *models.Token, error)
	UpdateToken(token string, expires time.Time) error
}

type tokenRepositoryGORM struct {
	db *gorm.DB
}

func NewTokenRepositoryGORM() TokenRepositoryGORM {
	return &tokenRepositoryGORM{
		db: config.InitDB(),
	}
}

func GenerateTokenPair(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["name"] = user.User_name
	claims["exp"] = time.Now().Add((time.Hour * 24) * 7).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetToken(c echo.Context) string {
	auth, exists := c.Request().Header["Authorization"]
	if !exists || len(auth) == 0 {
		return ""
	}

	Bearer := auth[0]
	tokenParts := strings.Split(Bearer, "Bearer ")
	if len(tokenParts) != 2 {
		return ""
	}

	token := tokenParts[1]
	return token
}

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"status":  "error",
				"message": "Unauthorized access",
			}
			return c.JSON(http.StatusUnauthorized, response)
		}

		token := parts[1]
		c.Set("token", token)
		return next(c)
	}
}

func (repo *tokenRepositoryGORM) SaveToken(user *models.User, token string) error {
	// Delete any existing token for the user
	if err := repo.db.Where("user_id = ?", user.User_id).Delete(&models.Token{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create and save the new token
	return repo.db.Create(&models.Token{
		UserId:    uuid.MustParse(user.User_id),
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}).Error
}

func (repo *tokenRepositoryGORM) SaveTokenResetEmail(user *models.User, token string) error {
	// Delete any existing token for the user
	if err := repo.db.Where("user_id = ?", user.User_id).Delete(&models.Token{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create and save the new token
	return repo.db.Create(&models.Token{
		UserId:    uuid.MustParse(user.User_id),
		Token:     token,
		ExpiresAt: time.Now().Add(time.Minute * 10),
	}).Error
}

func (repo *tokenRepositoryGORM) UpdateToken(token string, expires time.Time) error {
	return repo.db.Model(&models.Token{}).Where("token = ?", token).Update("expires_at", expires).Error
}

func (repo *tokenRepositoryGORM) GetUserIdByToken(token string) (string, error) {
	var AccessToken models.Token
	var User models.User

	err := repo.db.Where("token = ?", token).Where("expires_at > ?", time.Now()).Take(&AccessToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid or expired token")
		}
		return "", errors.New("database error while fetching token")
	}

	err = repo.db.Where("user_id = ?", AccessToken.UserId).First(&User).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "nil", errors.New("user not found")
		}
		return "", errors.New("database error while fetching user")
	}

	return User.User_id, nil

}

func (repo *tokenRepositoryGORM) GetUserByToken(token string) (*models.User, *models.Token, error) {
	var AccessToken models.Token
	var User models.User

	err := repo.db.Where("token = ?", token).Where("expires_at > ?", time.Now()).Take(&AccessToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid or expired token")
		}
		return nil, nil, errors.New("database error while fetching token")
	}

	err = repo.db.Where("user_id = ?", AccessToken.UserId).First(&User).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, errors.New("database error while fetching user")
	}

	return &User, &AccessToken, nil

}

func (repo *tokenRepositoryGORM) GetTokenByToken(token string) (*models.Token, error) {
	var t models.Token
	err := repo.db.Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}
