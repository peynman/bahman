package router

import (
	"errors"
	"fmt"
	"github.com/peyman-abdi/avalanche/app/interfaces/services"
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strings"
)

type routerImpl struct {
	server      *fasthttp.Server
	router      *routing.Router
	groups      map[string][]*services.RouteGroup
	middleWares map[string]*services.MiddleWare
	domain      string
	port        string
	log         services.Logger
	engine      services.RenderEngine
}

var _ services.Router = (*routerImpl)(nil)

func Initialize(app services.Application, config services.Config, logger services.Logger, engine services.RenderEngine) services.Router {
	s := new(routerImpl)

	s.router = routing.New()
	s.server = &fasthttp.Server{
		Handler: s.router.HandleRequest,
		Name:    config.GetAsString("server.name", "Avalanche"),
	}
	s.groups = make(map[string][]*services.RouteGroup)
	s.middleWares = make(map[string]*services.MiddleWare)

	s.router.NotFound(handleNotFound)
	s.log = logger
	s.engine = engine

	s.port = config.GetAsString("server.port", "8080")
	s.domain = config.GetAsString("server.address", "127.0.0.1")

	errTempPath := ""
	if app.IsDebugMode() {
		errTempPath = app.TemplatesPath(config.GetString("server.error", "services/error-debug.jet"))
	} else {
		errTempPath = app.TemplatesPath(config.GetString("server.error", "services/error-deploy.jet"))
	}
	engine.ParseTemplate(&services.Template{
		Name: "error",
		Path: errTempPath,
	})

	return s
}

func (r *routerImpl) Serve() error {
	return r.server.ListenAndServe(r.domain + ":" + r.port)
}
func (r *routerImpl) RegisterRoutes(routes []*services.Route) error {
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
						r.groups[group] = []*services.RouteGroup{}
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

func (r *routerImpl) RegisterMiddleWares(middleWares []*services.MiddleWare) error {
	for _, middleWare := range middleWares {
		if r.middleWares[middleWare.Name] != nil {
			return errors.New(fmt.Sprintf("Middleware with this name already exists: %s", middleWare.Name))
		}

		r.middleWares[middleWare.Name] = middleWare
	}
	return nil
}

func (r *routerImpl) RegisterGroups(groups []*services.RouteGroup) error {
	for _, group := range groups {
		if r.groups[group.Name] == nil {
			r.groups[group.Name] = []*services.RouteGroup{group}
		} else {
			r.groups[group.Name] = append(r.groups[group.Name], group)
		}
	}

	return nil
}

func (r *routerImpl) RemoveRoutes(routes []*services.Route) {
}
func (r *routerImpl) RemoveMiddleWares(middleWares []*services.MiddleWare) {
	for _, middleWare := range middleWares {
		delete(r.middleWares, middleWare.Name)
	}
}
func (r *routerImpl) RemoveGroups(groups []*services.RouteGroup) {
	for _, group := range groups {
		delete(r.groups, group.Name)
	}
}
