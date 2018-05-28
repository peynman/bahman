package router

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
	"strings"
	"errors"
	"fmt"
)

type routerImpl struct {
	server *fasthttp.Server
	router *routing.Router
	groups map[string][]interfaces.RequestHandler
	middleWares map[string]interfaces.RequestHandler
	domain string
	port string
}
var _ interfaces.Router = (*routerImpl)(nil)

func Initialize(config interfaces.Config) interfaces.Router {
	s := new(routerImpl)

	s.router = routing.New()
	s.server = &fasthttp.Server{
		Handler: s.router.HandleRequest,
		Name: config.GetAsString("server.name", "Avalanche"),
	}
	s.groups = make(map[string][]interfaces.RequestHandler)
	s.middleWares = make(map[string]interfaces.RequestHandler)

	s.router.NotFound(handleNotFound)

	s.port = config.GetAsString("server.port", "8080")
	s.domain = config.GetAsString("server.address", "127.0.0.1")

	return s
}

func LoadModules(modules []interfaces.Module) error {
	return nil
}

func (r *routerImpl) Serve() error {
	return r.server.ListenAndServe(r.domain + ":" + r.port)
}
func (r *routerImpl) RegisterRoutes(routes []*interfaces.Route) error {
	for _, route := range routes {
		var parent = &r.router.RouteGroup
		groups := strings.Split(route.Group, "/")
		if len(groups) > 0 {
			for _, group := range groups {
				if r.groups[group] == nil {
					r.groups[group] = []interfaces.RequestHandler {}
				}
				handlers := handleMiddleWares(r.groups[group])
				parent = parent.Group(group, handlers...)
			}
		}

		for _, middleWare := range route.MiddleWares {
			if r.middleWares[middleWare] != nil {
				parent.Use(handleMiddleWare(r.middleWares[middleWare]))
			}
		}

		methods := methodsFromInt(route.Methods)
		for _, method := range methods {
			parent.To(method, route.Url, handleRoute(route))
		}
	}
	return nil
}

func (r *routerImpl) RegisterMiddleWares(middleWares map[string]interfaces.RequestHandler) error  {
	for name, handler := range middleWares {
		if r.middleWares[name] != nil {
			return errors.New(fmt.Sprintf("Middleware with this name already exists: %s", name))
		}

		r.middleWares[name] = handler
	}
	return nil
}

func (r *routerImpl) RegisterGroups(groups map[string]interfaces.RequestHandler) error  {
	for name, handler := range groups {
		if r.groups[name] == nil {
			r.groups[name] = []interfaces.RequestHandler {handler}
		} else {
			r.groups[name] = append(r.groups[name], handler)
		}
	}

	return nil
}

func (r *routerImpl) RemoveRoutes(routes []*interfaces.Route) {
}
func (r *routerImpl) RemoveMiddleWares(middleWares map[string]interfaces.RequestHandler) {
}
func (r *routerImpl) RemoveGroups(groups map[string]interfaces.RequestHandler) {
}
