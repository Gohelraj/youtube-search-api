package error

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrInvalidValueInLimit  = generateError(http.StatusBadRequest, "invalid value in limit")
	ErrInvalidValueInOffset = generateError(http.StatusBadRequest, "invalid value in offset")
	ErrSearchStringRequired = generateError(http.StatusBadRequest, "searchString is required in request body")
	ErrLimitExceeded        = generateError(http.StatusBadRequest, "limit must be less than 100")
)

type Error struct {
	HttpStatusCode int    `json:"code,omitempty"`
	ErrMessage     string `json:"message"`
}

func (e Error) Error() string {
	return e.ErrMessage
}

func generateError(httpErrCode int, msg string) error {
	return &Error{
		HttpStatusCode: httpErrCode,
		ErrMessage:     msg,
	}
}

func SendError(c *gin.Context, err error) {
	if customErr, ok := err.(*Error); ok {
		c.JSON(customErr.HttpStatusCode, gin.H{"code": customErr.HttpStatusCode, "error": customErr.ErrMessage})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
	}
}
