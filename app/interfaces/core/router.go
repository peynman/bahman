package core

const (
	GET    = 1 << iota
	PUT    = 1 << iota
	POST   = 1 << iota
	DELETE = 1 << iota
	PATCH  = 1 << iota
	ANY    = GET | PUT | POST | DELETE | PATCH
)

type (
	RequestHandler func(Request, Response) error
)

type MiddleWare struct {
	Name     string
	Handler  RequestHandler
	Priority int
}

type RouteGroup struct {
	Name    string
	Handler RequestHandler
}

type Route struct {
	Name        string
	Description string
	Group       string
	MiddleWares []string
	Urls        []string
	Methods     int
	Verify      func(Request) error
	Handle      RequestHandler
}

type Request interface {
	SetValue(key string, val interface{})
	GetValue(key string) interface{}
	GetAll(names ...string) map[string]interface{}
}

type Response interface {
	SuccessString(content string) Response
	SuccessJSON(json interface{}) Response
	ContentType(contentType string) Response
	View(name string, params map[string]interface{}) Response
}

type Router interface {
	RegisterRoutes(routes []*Route) error
	RegisterMiddleWares(middleWares []*MiddleWare) error
	RegisterGroups(groups []*RouteGroup) error

	RemoveRoutes(routes []*Route)
	RemoveMiddleWares(middleWares []*MiddleWare)
	RemoveGroups(groups []*RouteGroup)

	Serve() error
}
