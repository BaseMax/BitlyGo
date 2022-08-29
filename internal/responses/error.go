package responses

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/itsjoniur/bitlygo/internal/durable"
	"github.com/unrolled/render"
)

// Struct for error responses
type ErrorResponse struct {
	Status  bool   `json:"status" default:"false"`
	Message string `json:"message"`
}

// BadRequestError means the request body is not a valid format
func BadRequestError(ctx context.Context, w http.ResponseWriter) {
	r := ctx.Value(2).(*render.Render)
	description := "The request body can not be parsed as valid data"
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusBadRequest, resp)
}

// NotFoundError means recourse is not found
func NotFoundError(ctx context.Context, w http.ResponseWriter) {
	r := ctx.Value(2).(*render.Render)
	description := http.StatusText(http.StatusNotFound)
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusNotFound, resp)
}

// InternalServerError means something not expected happend
func InternalServerError(ctx context.Context, w http.ResponseWriter) {
	r := ctx.Value(2).(*render.Render)
	description := http.StatusText(http.StatusInternalServerError)
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusInternalServerError, resp)
}

// InvalidLinkError means the link value is not a valid url
func InvalidLinkError(ctx context.Context, w http.ResponseWriter) {
	r := ctx.Value(2).(*render.Render)
	description := "link value must be a valid url"
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusBadRequest, resp)
}

// LinkIsExistsError means link is already exists in database
func LinkIsExistsError(ctx context.Context, w http.ResponseWriter, name string) {
	r := ctx.Value(2).(*render.Render)
	description := fmt.Sprintf("link with name %s exists", name)
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusConflict, resp)
}

// FieldEmptyError means a required field is empty
func FieldEmptyError(ctx context.Context, w http.ResponseWriter, field string) {
	r := ctx.Value(2).(*render.Render)
	description := fmt.Sprintf("field %s can not be empty", field)
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusBadRequest, resp)
}

// LimitRangeError means the limit value is out of expected range
func LimitRangeError(ctx context.Context, w http.ResponseWriter) {
	r := ctx.Value(2).(*render.Render)
	description := "limit value must be between 1-100"
	resp := createErr(ctx, description)

	r.JSON(w, http.StatusBadRequest, resp)
}

// CreateErr create a new ErrorResponse
func createErr(ctx context.Context, description string) *ErrorResponse {
	logger := ctx.Value(1).(*durable.Logger)

	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	trace := fmt.Sprintf("[ERROR] %s -> %s:%d %s", description, file, line, funcName)

	logger.Error(trace)

	return &ErrorResponse{
		Message: description,
	}
}
