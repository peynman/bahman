package hash_test

import (
	"testing"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/avest"
)

var s services.Services

func init() {
	s = avest.MockServices(avest.CommonConfigs, avest.CommonEnvs)
}

func TestNew(t *testing.T) {
	testStrings := []string {
		"",
		"a",
		"some string",
		"a really long string can go as a input too without problems",
	}
	for _, str := range testStrings {
		hash := s.Hash().Make(str)
		if !s.Hash().Compare(str, hash) {
			t.Errorf("Compared hashes are not same %s != %s", str, hash)
		}
	}
}
