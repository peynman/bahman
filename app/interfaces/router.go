package interfaces

const (
	GET = 1 << iota
	PUT = 1 << iota
	POST = 1 << iota
	DELETE = 1 << iota
	ANY = GET | PUT | POST | DELETE
)

type (
	RequestHandler func(Request, Response) error
)

type Route struct {
	Name 			string
	Description 	string
	Url				string
	Group			string
	MiddleWares 	[]string
	Methods			int
	Verify			func(Request) error
	Handle			RequestHandler
}

type Request interface {
	Param(key string) interface{}
}

type Response interface {
	Json() Response
}

type Router interface {
	RegisterRoutes(routes []*Route) error
	RegisterMiddleWares(middleWares map[string]RequestHandler) error
	RegisterGroups(groups map[string]RequestHandler) error

	RemoveRoutes(routes []*Route)
	RemoveMiddleWares(middleWares map[string]RequestHandler)
	RemoveGroups(groups map[string]RequestHandler)

	Serve() error
}
