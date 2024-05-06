package app

import (
	"context"
	"encoding/json"
	"golang.org/x/exp/slog"
	"net/http"
	"ps-cats-social/pkg/errs"
	"strings"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	hasBody         bool // For POST, PUT
	data            map[string]interface{}
	internalContext context.Context
}

func NewContext(rw http.ResponseWriter, req *http.Request) *Context {
	ctx := &Context{
		Writer:          rw,
		Request:         req,
		hasBody:         req.Method != http.MethodGet,
		internalContext: req.Context(),
	}

	if strings.Contains(ctx.Request.Header.Get("Content-type"), "application/json") {
		ctx.ParseJson()
	}

	return ctx
}

func (ctx *Context) GetJsonBody() map[string]interface{} {
	if ctx.hasBody {
		return ctx.data
	}
	return nil
}

func (ctx *Context) ParseJson() {
	if ctx.hasBody {
		decoder := json.NewDecoder(ctx.Request.Body)
		decoder.Decode(&ctx.data)
	} else {
		ctx.data = make(map[string]interface{})
	}
}

func (ctx *Context) Context() context.Context {
	if ctx.internalContext == nil {
		ctx.internalContext = ctx.Request.Context()
	}
	return ctx.internalContext
}

// for now default is error
func (ctx *Context) log(message string, level slog.Level) {
	slog.ErrorCtx(ctx.Context(), message, "attrs", errs.GetDefaultRequestFields(ctx.Request))
}

func (ctx *Context) ErrorLog(message string, level slog.Level) {
	ctx.log(message, level)
}
