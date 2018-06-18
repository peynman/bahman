package router

import (
	"encoding/json"
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
)

type responseImpl struct {
	context *routing.Context
	log     core.Logger
}

func (r *responseImpl) SuccessString(content string) core.Response {
	r.context.SuccessString("text/string", content)
	return r
}

func (r *responseImpl) SuccessJSON(object interface{}) core.Response {
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

func (r *responseImpl) ContentType(contentType string) core.Response {
	r.context.SetContentType(contentType)
	return r
}

var _ core.Response = (*responseImpl)(nil)
