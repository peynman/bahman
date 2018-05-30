package router

import (
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"encoding/json"
)

type responseImpl struct {
	context *routing.Context
	log interfaces.Logger
}

func (r *responseImpl) SuccessString(content string) interfaces.Response {
	r.context.SuccessString("text/string", content)
	return r
}

func (r *responseImpl) SuccessJSON(object interface{}) interfaces.Response {
	bytes, err := json.Marshal(object)
	if err != nil {
		r.log.FatalFields("Unable to serialize object to json", map[string]interface{} {
			"error": err,
			"object": object,
		})
		return r
	}

	r.context.Success("application/json", bytes)
	return r
}

func (r *responseImpl) ContentType(contentType string) interfaces.Response {
	r.context.SetContentType(contentType)
	return r
}

var _ interfaces.Response = (*responseImpl)(nil)

