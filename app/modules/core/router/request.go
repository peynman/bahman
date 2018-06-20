package router

import (
	"bytes"
	"encoding/json"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/peyman-abdi/fasthttp-routing"
)

type requestImpl struct {
	context        *routing.Context
	postJSONParams map[string]interface{}
}

var _ core.Request = (*requestImpl)(nil)

func (r *requestImpl) GetAll(names ...string) (result map[string]interface{}) {
	result = make(map[string]interface{})
	for _, name := range names {
		if r.context.Get(name) != nil {
			result[name] = r.context.Get(name)
		} else if r.context.QueryArgs().Has(name) {
			result[name] = r.context.QueryArgs().Peek(name)
		} else if r.context.PostArgs().Has(name) {
			result[name] = r.context.PostArgs().Peek(name)
		} else if r.postJSONParams != nil && r.postJSONParams[name] != nil {
			result[name] = r.postJSONParams[name]
		} else if r.context.HasParam(name) {
			result[name] = r.context.Param(name)
		}
	}
	return
}
func (r *requestImpl) GetValue(name string) interface{} {
	if r.context.Get(name) != nil {
		return r.context.Get(name)
	} else if r.context.QueryArgs().Has(name) {
		return r.context.QueryArgs().Peek(name)
	} else if r.context.PostArgs().Has(name) {
		return r.context.PostArgs().Peek(name)
	} else if r.postJSONParams != nil && r.postJSONParams[name] != nil {
		return r.postJSONParams[name]
	} else if r.context.HasParam(name) {
		return r.context.Param(name)
	}

	return nil
}
func (r *requestImpl) SetValue(name string, val interface{}) {
	r.context.Set(name, val)
}

func (r *routerImpl) requestFromContext(ctx *routing.Context) core.Request {
	var jsonParams map[string]interface{}

	if bytes.HasPrefix(ctx.Request.Header.ContentType(), []byte("application/json")) {
		err := json.Unmarshal(ctx.Request.Body(), &jsonParams)
		if err != nil {
			jsonParams = nil
		}
	}
	return &requestImpl{context: ctx, postJSONParams: jsonParams}
}
