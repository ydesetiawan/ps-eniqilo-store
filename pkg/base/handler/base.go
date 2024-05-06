package handler

import (
	"fmt"
	"go/types"
	"net/http"
	"ps-cats-social/pkg/base/app"
	"ps-cats-social/pkg/errs"
	"ps-cats-social/pkg/httphelper"
	"ps-cats-social/pkg/httphelper/response"
	"ps-cats-social/pkg/middleware"

	"golang.org/x/exp/slog"
)

type HandlerFn func(*app.Context) *response.WebResponse

type BaseHTTPHandler struct {
	Handlers interface{}
	DB       types.Nil
	Params   map[string]string
}

func (h *BaseHTTPHandler) RunAction(fn HandlerFn) http.HandlerFunc {
	return h.HandlePanic(h.Execute(fn))
}

func (h *BaseHTTPHandler) Execute(handler HandlerFn) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			// return error if err
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					httphelper.WriteJSON(rw, http.StatusInternalServerError,
						response.WebResponse{
							Status:  http.StatusInternalServerError,
							Message: http.StatusText(http.StatusInternalServerError),
						},
					)
					return
				}

				resultError := errs.ErrorAdvisor(err)
				httphelper.WriteJSON(rw, resultError.Status,
					response.WebResponse{
						Status:  resultError.Status,
						Message: resultError.Message,
						Error:   resultError.Error,
						Data:    types.Interface{},
					},
				)
			}
		}()

		ctx := app.NewContext(rw, r)
		resp := handler(ctx)
		httpStatus := resp.Status

		httphelper.WriteJSON(rw, httpStatus,
			response.WebResponse{
				Status:     httpStatus,
				Message:    resp.Message,
				Data:       resp.Data,
				Pagination: resp.Pagination})
	}
}

func (h *BaseHTTPHandler) HandlePanic(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.logPanicMessage(r, "CaptureLastPanic NEED TO FIX NOW", err)
				httphelper.WriteJSON(rw, http.StatusInternalServerError, "Request unhalted unxpectedly, Please contact administrator")
			}
		}()
		fn(rw, r)
	}
}

func (h *BaseHTTPHandler) logPanicMessage(r *http.Request, message string, err interface{}) {
	errStack := errs.StackAndFile(3)
	errInfo := fmt.Sprintf("\n SCM-production service \n* MUST FIX \U0001f4a3 \U0001f4a3 \U0001f4a3 "+
		"Panic Error: %v*", err)
	msg := fmt.Sprintf("%s\n\nStack trace: \n%s...", errInfo, errStack)

	fmt.Println("\nPANIC:", msg)
	src := "\n--- (Staging " + r.Host + ") ---\n"

	slog.ErrorCtx(r.Context(), message+src+msg, "attrs", errs.GetDefaultRequestFields(r))
}

func (h *BaseHTTPHandler) RunActionAuth(fn HandlerFn) http.HandlerFunc {
	return h.HandlePanic(h.ExecuteAuth(h.Execute(fn)))
}

func (h *BaseHTTPHandler) ExecuteAuth(fn http.HandlerFunc) http.HandlerFunc {
	return middleware.JWTAuthMiddleware(fn)
}
