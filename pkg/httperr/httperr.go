package httperr

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	wctx "github.com/oakeshq/go-starter/context"

	"github.com/oakeshq/go-starter/pkg/logs"
	"github.com/oakeshq/go-starter/pkg/render"
)

var oauthErrorMap = map[int]string{
	http.StatusBadRequest:          "invalid_request",
	http.StatusUnauthorized:        "unauthorized_client",
	http.StatusForbidden:           "access_denied",
	http.StatusInternalServerError: "server_error",
	http.StatusServiceUnavailable:  "temporarily_unavailable",
}

// OAuthError is the JSON handler for OAuth2 error responses
type OAuthError struct {
	Err             string `json:"error"`
	Description     string `json:"error_description,omitempty"`
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
}

func (e *OAuthError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%s: %s", e.Err, e.Description)
}

// WithInternalError adds internal error information to the error
func (e *OAuthError) WithInternalError(err error) *OAuthError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *OAuthError) WithInternalMessage(fmtString string, args ...interface{}) *OAuthError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

// Cause returns the root cause error
func (e *OAuthError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

func OauthError(err string, description string) *OAuthError {
	return &OAuthError{Err: err, Description: description}
}

func BadRequestError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, fmtString, args...)
}

func InternalServerError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, fmtString, args...)
}

func NotFoundError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusNotFound, fmtString, args...)
}

func UnauthorizedError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, fmtString, args...)
}

func ForbiddenError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusForbidden, fmtString, args...)
}

func UnprocessableEntityError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, fmtString, args...)
}

// HTTPError is an error with a message and an HTTP status code.
type HTTPError struct {
	Code            int    `json:"code"`
	Message         string `json:"msg"`
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
	ErrorID         string `json:"error_id,omitempty"`
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Cause returns the root cause error
func (e *HTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

// WithInternalError adds internal error information to the error
func (e *HTTPError) WithInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

func httpError(code int, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(fmtString, args...),
	}
}

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
func Recoverer(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	defer func() {
		if rvr := recover(); rvr != nil {

			logEntry := logs.GetLogEntry(r)
			if logEntry != nil {
				logEntry.Panic(rvr, debug.Stack())
			} else {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()
			}

			err := &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
			HandleError(err, w, r)
		}
	}()

	return nil, nil
}

type ErrorCause interface {
	Cause() error
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	log := logs.GetLogEntry(r)
	errorID := wctx.GetRequestID(r.Context())
	switch e := err.(type) {
	case *HTTPError:
		if e.Code >= http.StatusInternalServerError {
			e.ErrorID = errorID
			// this will get us the stack trace too
			log.WithError(e.Cause()).Error(e.Error())
		} else {
			log.WithError(e.Cause()).Warn(e.Error())
		}
		if jsonErr := render.SendJSON(w, e.Code, e); jsonErr != nil {
			HandleError(jsonErr, w, r)
		}
	case *OAuthError:
		log.WithError(e.Cause()).Info(e.Error())
		if jsonErr := render.SendJSON(w, http.StatusBadRequest, e); jsonErr != nil {
			HandleError(jsonErr, w, r)
		}
	case ErrorCause:
		HandleError(e.Cause(), w, r)
	default:
		log.WithError(e).Sugar().Errorf("Unhandled server error: %s", e.Error())
		// hide real error details from response to prevent info leaks
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(`{"code":500,"msg":"Internal server error","error_id":"` + errorID + `"}`)); writeErr != nil {
			log.WithError(writeErr).Error("Error writing generic error message")
		}
	}
}
