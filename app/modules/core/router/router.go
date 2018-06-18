package router

import (
	"errors"
	"fmt"
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strings"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
)

type routerImpl struct {
	server      *fasthttp.Server
	router      *routing.Router
	groups      map[string][]*core.RouteGroup
	middleWares map[string]*core.MiddleWare
	domain      string
	port        string
	log         core.Logger
}

var _ core.Router = (*routerImpl)(nil)

func Initialize(config core.Config, logger core.Logger) core.Router {
	s := new(routerImpl)

	s.router = routing.New()
	s.server = &fasthttp.Server{
		Handler: s.router.HandleRequest,
		Name:    config.GetAsString("server.name", "Avalanche"),
	}
	s.groups = make(map[string][]*core.RouteGroup)
	s.middleWares = make(map[string]*core.MiddleWare)

	s.router.NotFound(handleNotFound)
	s.log = logger

	s.port = config.GetAsString("server.port", "8080")
	s.domain = config.GetAsString("server.address", "127.0.0.1")

	return s
}

func (r *routerImpl) Serve() error {
	return r.server.ListenAndServe(r.domain + ":" + r.port)
}
func (r *routerImpl) RegisterRoutes(routes []*core.Route) error {
	for _, route := range routes {
		var parent = &r.router.RouteGroup
		if len(route.Group) > 0 {
			groups := strings.Split(route.Group, "/")
			if len(groups) > 0 {
				for _, group := range groups {
					if len(group) == 0 {
						continue
					}

					if r.groups[group] == nil {
						r.groups[group] = []*core.RouteGroup{}
					}
					handlers := r.handleGroups(r.groups[group])
					parent = parent.Group("/" + group)
					parent.Use(handlers...)
				}
			}
		}

		for _, middleWare := range route.MiddleWares {
			if r.middleWares[middleWare] != nil {
				parent.Use(r.handleMiddleWare(r.middleWares[middleWare].Handler))
			}
		}

		methods := r.methodsFromInt(route.Methods)
		for _, method := range methods {
			for _, url := range route.Urls {
				parent.To(method, url, r.handleRoute(route)).Name(route.Name)
			}
		}
	}
	return nil
}

func (r *routerImpl) RegisterMiddleWares(middleWares []*core.MiddleWare) error {
	for _, middleWare := range middleWares {
		if r.middleWares[middleWare.Name] != nil {
			return errors.New(fmt.Sprintf("Middleware with this name already exists: %s", middleWare.Name))
		}

		r.middleWares[middleWare.Name] = middleWare
	}
	return nil
}

func (r *routerImpl) RegisterGroups(groups []*core.RouteGroup) error {
	for _, group := range groups {
		if r.groups[group.Name] == nil {
			r.groups[group.Name] = []*core.RouteGroup{group}
		} else {
			r.groups[group.Name] = append(r.groups[group.Name], group)
		}
	}

	return nil
}

func (r *routerImpl) RemoveRoutes(routes []*core.Route) {
}
func (r *routerImpl) RemoveMiddleWares(middleWares []*core.MiddleWare) {
	for _, middleWare := range middleWares {
		delete(r.middleWares, middleWare.Name)
	}
}
func (r *routerImpl) RemoveGroups(groups []*core.RouteGroup) {
	for _, group := range groups {
		delete(r.groups, group.Name)
	}
}
