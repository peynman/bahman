package hash

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"strings"
	"crypto/sha256"
	"hash"
	"crypto/sha512"
	"encoding/hex"
)

type hashImpl struct {
	hashier hash.Hash
}

func (a *hashImpl) Make(str string) string {
	a.hashier.Reset()
	_, err := a.hashier.Write([]byte(str))
	if err != nil {
		return ""
	}

	return hex.EncodeToString(a.hashier.Sum(nil))
}

func (a *hashImpl) Compare(bare string, hash string) bool {
	return a.Make(bare) == hash
}

func New(config services.Config) services.Hash {
	hash := &hashImpl{}

	algorithm := strings.ToLower(config.GetString("hash.algorithm", "sha256"))
	switch algorithm {
	case "sha256":
		hash.hashier = sha256.New()
	case "sha512":
		hash.hashier = sha512.New()
	default:
		hash.hashier = sha512.New()
	}
	return hash
}
