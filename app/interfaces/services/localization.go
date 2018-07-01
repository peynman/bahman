package services

type Localization interface {
	L(string) string
	LP(string, map[string]string) string
}
