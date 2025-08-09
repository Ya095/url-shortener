package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Response - часть response повторяющаяся для всех handlers
type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}

const (
	StatusOK = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: "OK",
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error: msg,
	}
}

// ValidationError - функция для корректного отображения ошибки при валидации
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required": 
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valied URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs, ", "),
	}
}
