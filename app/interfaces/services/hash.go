package services

type Hash interface {
	Make(string) string
	Compare(bare string, hashed string) bool
}