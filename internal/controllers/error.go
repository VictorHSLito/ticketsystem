package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ApiError struct {
	Status        int            `json:"-"`
	Title         string         `json:"title"`
	Type          string         `json:"type"`
	Detail        string         `json:"detail"`
	Instance      string         `json:"instance,omitempty"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Detail
}

var (
	ErrUserNotFound       = &ApiError{Status: http.StatusNotFound, Type: "about:blank", Title: "Not Found", Detail: "User not found!"}
	ErrEventNotFound      = &ApiError{Status: http.StatusNotFound, Type: "about:blank", Title: "Not Found", Detail: "Event not found!"}
	ErrInternalServer     = &ApiError{Status: http.StatusInternalServerError, Type: "about:blank", Title: "Internal Server Error", Detail: "An unexpected error occurred"}
	ErrDuplicatedResource = &ApiError{Status: http.StatusBadRequest, Type: "about:blank", Title: "Bad Request", Detail: "Duplicated resource!"}
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
				"type":           apiError.Type,
				"title":          apiError.Title,
				"status":         apiError.Status,
				"detail":         apiError.Detail,
				"instance":       c.Request.RequestURI,
				"invalid_params": apiError.InvalidParams,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"type":     "about:blank",
			"title":    ErrInternalServer.Title,
			"status":   http.StatusInternalServerError,
			"detail":   ErrInternalServer.Detail,
			"instance": c.Request.RequestURI,
		})
	}
}

func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			var invalidParams []InvalidParam

			for _, e := range validationErrs {
				invalidParams = append(invalidParams, InvalidParam{
					Name:   e.Field(),
					Reason: getReasonField(e),
				})
			}

			c.Error(&ApiError{
				Status:        http.StatusBadRequest,
				Title:         "Validation Error",
				Type:          "about:blank",
				Detail:        "The request contain invalid fields",
				Instance:      c.Request.RequestURI,
				InvalidParams: invalidParams,
			})

			return false
		}

		c.Error(&ApiError{
			Status: http.StatusBadRequest,
			Type:   "about:blank",
			Title:  "Malformed JSON",
			Detail: err.Error(), Instance: c.Request.RequestURI,
		})

	}
	return true
}

func getReasonField(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "gte":
		return "The value should be greater"
	default:
		return "Failed on rule: " + e.Tag()
	}
}
