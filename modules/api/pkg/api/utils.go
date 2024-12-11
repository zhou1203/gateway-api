package api

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
)

// Avoid emitting errors that look like valid HTML. Quotes are okay.
var sanitizer = strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")

func HandleInternalError(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusInternalServerError, response, req, err)
}

// HandleBadRequest writes http.StatusBadRequest and log error
func HandleBadRequest(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusBadRequest, response, req, err)
}

func HandleNotFound(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusNotFound, response, req, err)
}

func HandleForbidden(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusForbidden, response, req, err)
}

func HandleUnauthorized(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusUnauthorized, response, req, err)
}

func HandleTooManyRequests(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusTooManyRequests, response, req, err)
}

func HandleConflict(response http.ResponseWriter, req *http.Request, err error) {
	handle(http.StatusConflict, response, req, err)
}

func HandleError(response http.ResponseWriter, req *http.Request, err error) {
	var statusCode int
	switch t := err.(type) {
	case errors.APIStatus:
		statusCode = int(t.Status().Code)
	case restful.ServiceError:
		statusCode = t.Code
	default:
		statusCode = http.StatusInternalServerError
	}
	handle(statusCode, response, req, err)
}

func handle(statusCode int, response http.ResponseWriter, req *http.Request, err error) {
	_, fn, line, _ := runtime.Caller(2)
	klog.Errorf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), statusCode)
}
