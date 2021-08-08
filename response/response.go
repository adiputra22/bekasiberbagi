package response

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

const STATUS_SUCCESS = "success"
const STATUS_FAILED = "failed"

func APIResponseSuccess(message string, code int, data interface{}) Response {
	meta := Meta{}
	meta.Message = message
	meta.Code = code
	meta.Status = STATUS_SUCCESS

	response := Response{}
	response.Meta = meta
	response.Data = data

	return response
}

func APIResponseFailed(message string, code int) Response {
	meta := Meta{}
	meta.Message = message
	meta.Code = code
	meta.Status = STATUS_FAILED

	response := Response{}
	response.Meta = meta
	response.Data = nil

	return response
}

func APIResponseFailedWithData(message string, code int, data interface{}) Response {
	meta := Meta{}
	meta.Message = message
	meta.Code = code
	meta.Status = STATUS_FAILED

	response := Response{}
	response.Meta = meta
	response.Data = data

	return response
}

func APIResponseValidationFailed(message string, code int, errors error) Response {
	meta := Meta{}
	meta.Message = message
	meta.Code = code
	meta.Status = STATUS_FAILED

	response := Response{}
	response.Meta = meta

	formatError := formatValidationError(errors)
	errorData := gin.H{"errors": formatError}

	response.Data = errorData

	return response
}

func formatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
