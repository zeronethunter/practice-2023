package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"practice/internal/pkg/users"
)

func NewBasicAuth(users []users.User) func(string, string, echo.Context) (bool, error) {
	return func(username, password string, c echo.Context) (bool, error) {
		for _, user := range users {
			if username == user.Username {
				if password == user.Password {
					return true, nil
				}

				return false, errors.Errorf("wrong password for user: %s", username)
			}
		}

		return false, errors.Errorf("user %s not found", username)
	}
}
