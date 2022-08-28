package responses

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type ErrorResponse struct {
	Status  bool   `json:"status" default:"false"`
	Message string `json:"message"`
}

func BadRequestError(ctx context.Context, w http.ResponseWriter) {
	description := "The request body can not be parsed as valid data"
	createErr(ctx, description)
}

func NotFoundError(ctx context.Context, w http.ResponseWriter) {
	description := http.StatusText(http.StatusNotFound)
	createErr(ctx, description)
}

func InternalServerError(ctx context.Context, w http.ResponseWriter) {
	description := http.StatusText(http.StatusInternalServerError)
	createErr(ctx, description)
}

func InvalidLinkError(ctx context.Context, w http.ResponseWriter) {
	description := "link value must be a valid url"
	createErr(ctx, description)
}

func LinkIsExistsError(ctx context.Context, w http.ResponseWriter, name string) {
	description := fmt.Sprintf("link with name %s exists", name)
	createErr(ctx, description)
}

func FieldEmptyError(ctx context.Context, w http.ResponseWriter, field string) {
	description := fmt.Sprintf("field %s can not be empty", field)
	createErr(ctx, description)
}

func LimitRangeError(ctx context.Context, w http.ResponseWriter) {
	description := "limit value must be between 1-100"
	createErr(ctx, description)
}

func createErr(ctx context.Context, description string) *ErrorResponse {
	logger := ctx.Value(1).(*log.Logger)

	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	trace := fmt.Sprintf("[ERROR] %s -> %s:%d %s", description, file, line, funcName)

	logger.Println(trace)

	return &ErrorResponse{
		Message: description,
	}
}
