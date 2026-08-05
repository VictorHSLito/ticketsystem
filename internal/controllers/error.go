package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Status  int    `json:"-"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound       = &ApiError{http.StatusNotFound, "NOT_FOUND", "User not found!"}
	ErrEventNotFound      = &ApiError{http.StatusNotFound, "EVENT_NOT_FOUND", "Event not found!"}
	ErrInternalServer     = &ApiError{http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred"}
	ErrBindingResources   = &ApiError{http.StatusBadRequest, "BAD_REQUEST", "Some fields are missing or in wrong format"}
	ErrDuplicatedResource = &ApiError{http.StatusBadRequest, "BAD_REQUEST", "Duplicated resource!"}
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if apiError, ok := errors.AsType[*ApiError](err); ok {
			c.JSON(apiError.Status, gin.H{
				"title":   apiError.Title,
				"message": apiError.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":   ErrInternalServer.Title,
				"message": ErrInternalServer.Message,
			})
		}
	}
}
