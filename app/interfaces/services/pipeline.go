package services


type Request interface {
	SetValue(key string, val interface{})
	GetValue(key string) interface{}
	GetAll(names ...string) map[string]interface{}
	GetAllAsString(names ...string) map[string]string
	GetBody() []byte
	GetBodyAsMap() map[string]string
	GetMethod() int
	GetUrl() string

	Session() Session
}

type Response interface {
	Error(applicationError ApplicationError) Response
	SuccessString(content string) Response
	SuccessJSON(json interface{}) Response
	View(name string, params map[string]interface{}) error
	Redirect(url string, statusCode int) error
	RedirectToRouteName(name string, statusCode int) error

	ContentType(contentType string) Response
	StatusCode(statusCode int) Response
	AddCookie(name string, value string) Response
	RemoveCookie(name string) Response
}

type Session interface {
	GetID() string
	Get(key string) interface{}
	GetAsString(key string) string
	GetAsArray(key string) []interface{}
	GetAsMap(key string) map[string]interface{}
	GetAsInt(key string) int
	Exists(key string) bool
	Set(key string, val interface{})
	GetAll() map[string]interface{}
	Delete(key string)
	Clear()
}

