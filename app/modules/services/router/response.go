package router

import (
	"encoding/json"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"github.com/peyman-abdi/fasthttp-routing"
)

type responseImpl struct {
	context *routing.Context
	log     services.Logger
	engine  services.RenderEngine
}

func (r *responseImpl) SuccessString(content string) services.Response {
	r.context.SuccessString("text/string", content)
	return r
}

func (r *responseImpl) SuccessJSON(object interface{}) services.Response {
	bytes, err := json.Marshal(object)
	if err != nil {
		r.log.FatalFields("Unable to serialize object to json", map[string]interface{}{
			"error":  err,
			"object": object,
		})
		return r
	}

	r.context.Success("application/json", bytes)
	return r
}

func (r *responseImpl) ContentType(contentType string) services.Response {
	r.context.SetContentType(contentType)
	return r
}

func (r *responseImpl) StatusCode(statusCode int) services.Response {
	r.context.SetStatusCode(statusCode)
	return r
}

func (r *responseImpl) View(name string, params map[string]interface{}) services.Response {
	err := r.engine.Render(name, params, r.context.Response.BodyWriter())
	if err != nil {
		r.ContentType("text/html")
		r.StatusCode(500)
		r.engine.Render("error", nil, r.context.Response.BodyWriter())
		return r
	}

	r.ContentType("text/html")
	r.StatusCode(200)
	return r
}

var _ services.Response = (*responseImpl)(nil)
