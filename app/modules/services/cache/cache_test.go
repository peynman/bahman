package cache_test

import (
	"testing"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/avest"
)

var s services.Services

func init() {
	s = avest.MockServices(avest.CommonConfigs, avest.CommonEnvs)
}
func TestCache(t *testing.T) {
	s.Cache().Set("user:123", "{val:100}", nil)
	val := s.Cache().Get("user:123")
	if _, ok := val.(string); val == nil || !ok {
		t.Error("not valid string")
	}
	if value := val.(string); value != "{val:100}" {
		t.Errorf("values are not the same!")
	}
}