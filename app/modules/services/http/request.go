package http

import (
	"bytes"
	"encoding/json"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/fasthttp-routing"
	"strings"
)

type requestImpl struct {
	context        *routing.Context
	postJSONParams map[string]interface{}
	session services.Session
}

func (r *requestImpl) GetUrl() string {
	return r.context.URI().String()
}

func (r *requestImpl) GetAllAsString(names ...string) (result map[string]string) {
	result = make(map[string]string)
	for _, name := range names {
		if r.context.Get(name) != nil {
			result[name] = r.context.Get(name).(string)
		} else if r.context.QueryArgs().Has(name) {
			result[name] = string(r.context.QueryArgs().Peek(name))
		} else if r.context.PostArgs().Has(name) {
			result[name] = string(r.context.PostArgs().Peek(name))
		} else if r.postJSONParams != nil && r.postJSONParams[name] != nil {
			result[name] = r.postJSONParams[name].(string)
		} else if r.context.HasParam(name) {
			result[name] = r.context.Param(name)
		}
	}
	return
}


func (r *requestImpl) GetBody() []byte {
	return r.context.PostBody()
}

func (r *requestImpl) GetBodyAsMap() map[string]string {
	result := map[string]string {}
	args := r.context.PostArgs()
	args.VisitAll(func(key, value []byte) {
		result[string(key)] = string(value)
	})
	return result
}

func (r *requestImpl) GetMethod() int {
	switch strings.ToLower(string(r.context.Method())) {
	case "get": return services.GET
	case "post": return services.POST
	case "head": return services.HEAD
	case "patch": return services.PATCH
	case "put": return services.PUT
	case "delete": return services.DELETE
	}
	return 0
}

var _ services.Request = (*requestImpl)(nil)

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
func (r *requestImpl) Session() services.Session {
	return r.session
}

func (r *routerImpl) requestFromContext(ctx *routing.Context) services.Request {
	reqRef := ctx.Get("bahman_request_reference")
	if reqRef != nil {
		if req, ok := reqRef.(services.Request); ok {
			return req
		}
	}

	var jsonParams map[string]interface{}

	if bytes.HasPrefix(ctx.Request.Header.ContentType(), []byte("application/json")) {
		err := json.Unmarshal(ctx.Request.Body(), &jsonParams)
		if err != nil {
			jsonParams = nil
		}
	}

	var sess services.Session
	if r.autoSession {
		sess = NewSession(ctx)
	}

	req := &requestImpl{context: ctx, postJSONParams: jsonParams, session: sess}
	ctx.Set("bahman_request_reference", req)
	return req
}
