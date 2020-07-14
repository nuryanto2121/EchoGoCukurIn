package tool

import (
	"net/http"
	"nuryanto2121/dynamic_rest_api_go/models"

	"github.com/sirupsen/logrus"
)

// GetStatusCode :
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	case models.Unauthorized:
		return http.StatusUnauthorized
	case models.ErrInvalidLogin:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
