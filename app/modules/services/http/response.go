package http

import (
	"encoding/json"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/fasthttp-contrib/sessions"
	"github.com/valyala/fasthttp"
)

type responseImpl struct {
	context *routing.Context
	log     services.Logger
	engine  services.RenderEngine
	router services.Router
}

func (r *responseImpl) RedirectToRouteName(name string, statusCode int) error {
	r.context.Redirect(r.router.Url(name), statusCode)
	return nil
}

func (r *responseImpl) Redirect(url string, statusCode int) error {
	r.context.Redirect(url, statusCode)
	return nil
}

func (r *responseImpl) Error(applicationError services.ApplicationError) services.Response {
	return r
}

func (r *responseImpl) AddCookie(name string, value string) services.Response {
	cookie := fasthttp.AcquireCookie()
	cookie.SetValue(value)
	cookie.SetKey(name)
	sessions.AddFasthttpCookie(cookie, r.context.RequestCtx)
	return r
}

func (r *responseImpl) RemoveCookie(name string) services.Response {
	sessions.RemoveFasthttpCookie(name, r.context.RequestCtx)
	return r
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

func (r *responseImpl) View(name string, params map[string]interface{}) error {
	err := r.engine.Render(name, params, r.context.Response.BodyWriter())
	if err != nil {
		return err
	}

	r.ContentType("text/html")
	r.StatusCode(200)
	return nil
}

var _ services.Response = (*responseImpl)(nil)

func (r *routerImpl) responseFromContext(ctx *routing.Context) services.Response {
	resRef := ctx.Get("bahman_response_reference")
	if resRef != nil {
		if res, ok := resRef.(services.Response); ok {
			return res
		}
	}

	res := &responseImpl{context: ctx, log: r.log, engine: r.engine, router: r}
	ctx.Set("bahman_response_reference", res)
	return res
}
