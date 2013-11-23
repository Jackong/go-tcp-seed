/**
 * User: jackong
 * Date: 11/11/13
 * Time: 7:29 PM
 */
package global

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Jackong/go-web-seed/err"
	"github.com/Jackong/go-web-seed/io/output"
)

type BeforeFunc func(http.ResponseWriter, *http.Request) error

var (
	beforeFuncs []BeforeFunc
	afterFuncs []http.HandlerFunc
)

type router struct {
	*mux.Router
}

func (this *router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	op := output.NewJson(writer)

	defer func() {
		status := http.StatusOK
		msg := "ok"
		if e := recover(); e != nil {
			switch e.(type) {
			case err.AccessError:
				accessErr := e.(err.AccessError)
				Access.Error(accessErr.Status, req, accessErr.Msg)
				status = accessErr.Status
				msg = accessErr.Msg
			default:
				Access.Alert(-2, req, "unkown panic", e)
				Log.Alert("unkown panic", e)
				status = http.StatusInternalServerError
				msg = "internal server error"
			}
		}

		if re := op.Render(status, msg); re != nil {
			Access.Alert(-1, req, "render", re)
		}
	}()

	if !fireBefore(op, req, beforeFuncs) {
		return
	}
	this.Router.ServeHTTP(op, req)

	fireAfter(op, req)
}

func OnBefore(before ...BeforeFunc) {
    beforeFuncs = append(beforeFuncs, before...)
}

func fireBefore(writer http.ResponseWriter, req *http.Request, beforeFuncs []BeforeFunc) bool {
	for _, beforeFunc := range beforeFuncs {
		if accessErr := beforeFunc(writer, req); accessErr != nil {
			panic(accessErr)
			return false
		}
	}
	return true
}

func OnAfter(after ...http.HandlerFunc) {
    afterFuncs = append(afterFuncs, after...)
}

func fireAfter(writer http.ResponseWriter, req *http.Request) {
    for _, afterFunc := range afterFuncs {
		afterFunc(writer, req)
	}
}

func OnBeforeRouteFunc(handlerFunc http.HandlerFunc, before ...BeforeFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		if fireBefore(writer, req, before) {
			handlerFunc(writer, req)
		}
	}
}

func OnBeforeRoute(handler http.Handler, before ...BeforeFunc) http.Handler {
	return OnBeforeRouteFunc(handler.ServeHTTP, before...)
}
