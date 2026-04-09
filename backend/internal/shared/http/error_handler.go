package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FortiBrine/VoidShift/internal/shared"
	"github.com/FortiBrine/VoidShift/internal/shared/http/validator"
	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func CustomErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}

	if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
		err := c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "validation_error",
			Message: "validation failed",
			Errors:  validationErr.Fields,
		})

		if err != nil {
			fmt.Printf("failed to send validation error response: %v\n", err)
		}
		return
	}

	if httpErr, ok := errors.AsType[*echo.HTTPError](err); ok {
		msg := httpErr.Message

		_ = c.JSON(httpErr.Code, ErrorResponse{
			Code:    "http_error",
			Message: msg,
		})
		return
	}

	if appErr, ok := errors.AsType[*shared.AppError](err); ok {
		statusCode := appErr.StatusCode
		code := appErr.Code
		msg := appErr.Message

		_ = c.JSON(statusCode, ErrorResponse{
			Code:    code,
			Message: msg,
		})
		return
	}

	code := echo.StatusCode(err)
	if code != 0 {
		_ = c.JSON(code, ErrorResponse{
			Code:    "http_error",
			Message: http.StatusText(code),
		})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "internal_error",
		Message: "internal server error",
	})
}
