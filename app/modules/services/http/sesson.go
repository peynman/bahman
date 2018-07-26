package http

import (
	"github.com/fasthttp-contrib/sessions"
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
)

type sessionImpl struct {
	session sessions.Session
}

func (s *sessionImpl) GetID() string {
	return s.session.ID()
}

func (s *sessionImpl) Get(key string) interface{} {
	return s.session.Get(key)
}

func (s *sessionImpl) GetAsString(key string) string {
	return s.session.GetString(key)
}

func (s *sessionImpl) GetAsInt(key string) int {
	return s.session.GetInt(key)
}

func (s *sessionImpl) GetAsArray(key string) []interface{} {
	if val, ok := s.session.Get(key).([]interface{}); ok {
		return val
	}

	return nil
}

func (s *sessionImpl) GetAsMap(key string) map[string]interface{} {
	if val, ok := s.session.Get(key).(map[string]interface{}); ok {
		return val
	}

	return nil
}

func (s *sessionImpl) Exists(key string) bool {
	return s.session.Get(key) != nil
}

func (s *sessionImpl) Set(key string, val interface{}) {
	s.session.Set(key, val)
}

func (s *sessionImpl) GetAll() map[string]interface{} {
	return s.session.GetAll()
}

func (s *sessionImpl) Delete(key string) {
	s.session.Delete(key)
}

func (s *sessionImpl) Clear() {
	s.session.Clear()
}

func NewSession(ctx *routing.Context) services.Session {
	sess := sessions.StartFasthttp(ctx.RequestCtx)
	return &sessionImpl{
		session: sess,
	}
}
