package services

const (
	GET    = 1 << iota
	PUT    = 1 << iota
	POST   = 1 << iota
	DELETE = 1 << iota
	PATCH  = 1 << iota
	HEAD   = 1 << iota
	ANY    = GET | PUT | POST | DELETE | PATCH | HEAD
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

type HttpError interface {
	ApplicationError
	StatusCode() int
	Request() Request
	Response() Response
}

type Router interface {
	RegisterRoutes(routes []*Route) error
	RegisterMiddleWares(middleWares []*MiddleWare) error
	RegisterGroups(groups []*RouteGroup) error

	RemoveRoutes(routes []*Route)
	RemoveMiddleWares(middleWares []*MiddleWare)
	RemoveGroups(groups []*RouteGroup)

	Load(services Services) error
	Url(name string) string
	UseGlobally([]*MiddleWare)
	Serve() error
	ServeTSL() error
}
