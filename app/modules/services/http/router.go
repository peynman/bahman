package http

import (
	"errors"
	"fmt"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strings"
	"github.com/fasthttp-contrib/sessions"
	"time"
	"github.com/peyman-abdi/bahman/app/modules/services/http/dialects"
)

type routerImpl struct {
	server      *fasthttp.Server
	fs   		*fasthttp.FS
	router      *routing.Router
	fsHandler   fasthttp.RequestHandler
	groups      map[string][]*services.RouteGroup
	middleWares map[string]*services.MiddleWare
	domain      string
	port        string
	cert   		string
	key string
	log         services.Logger
	engine      services.RenderEngine
	autoSession bool
}

var _ services.Router = (*routerImpl)(nil)

func (r *routerImpl) Serve() error {
	return r.server.ListenAndServe(r.domain + ":" + r.port)
}

func (r *routerImpl) ServeTSL() error {
	return r.server.ListenAndServeTLS(r.domain + ":" + r.port, r.cert, r.key)
}

func (r *routerImpl) Url(name string)string {
	if route := r.router.Route(name); route != nil {
		return route.URL()
	}
	return ""
}

func (r *routerImpl) RegisterRoutes(routes []*services.Route) error {
	for _, route := range routes {
		var parent = &r.router.RouteGroup
		var finalHandles []routing.Handler
		for _, middleWare := range route.MiddleWares {
			if r.middleWares[middleWare] != nil {
				finalHandles = append(finalHandles, r.handleMiddleWare(r.middleWares[middleWare].Handler))
			}
		}
		setMiddleWares := false

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
					if !setMiddleWares {
						setMiddleWares = true
						finalHandles = append(finalHandles, handlers...)
						parent.Use(finalHandles...)
					} else {
						parent.Use(handlers...)
					}
				}
			}
		}


		methods := r.methodsFromInt(route.Methods)
		for _, method := range methods {
			for _, url := range route.Urls {
				if !setMiddleWares {
					setMiddleWares = true
					finalHandles = append(finalHandles, r.handleRoute(route))
					parent.To(method, url, finalHandles...).Name(route.Name)
				} else {
					parent.To(method, url, r.handleRoute(route)).Name(route.Name)
				}
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

func (r *routerImpl) UseGlobally(wares []*services.MiddleWare)  {
	r.router.Use(r.handleMiddleWares(wares)...)
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
func (s *routerImpl)  Load(instance services.Services) error {
	config := instance.Config()
	engine := instance.Renderer()

	if config.GetBoolean("server.router.error.enabled", true) {
		engine.ParseTemplate(&services.Template {
			Name: config.GetString("server.router.error.name", "http.error"),
			Path: config.GetString("server.router.error.template", "http/error.jet"),
		})
	}
	if config.GetBoolean("server.router.reports.enabled", true) {
		engine.ParseTemplate(&services.Template {
			Name: config.GetString("server.router.reports.name", "http.routes"),
			Path: config.GetString("server.router.reports.template", "http/router.jet"),
		})
		s.RegisterRoutes([]*services.Route {
			{
				Name: config.GetString("server.router.reports.name", "http.routes"),
				Methods: services.GET,
				MiddleWares: []string{},
				Group: config.GetString("server.router.reports.group", "router"),
				Urls: []string{config.GetString("server.router.reports.url", "")},
				Handle: func(request services.Request, response services.Response) error {
					return response.View(config.GetString("server.router.reports.name", "http.routes"), map[string]interface{} {
						"route": map[string]string {
						"url": "something",
					}})
				},
			},
		})
	}
	if config.GetBoolean("server.router.home.enabled", true) {
		engine.ParseTemplate(&services.Template{
			Name: config.GetString("server.router.home.name", "http.home"),
			Path: config.GetString("server.router.home.path", "http/home.jet"),
		})
		s.RegisterRoutes([]*services.Route {
			{
				Name: config.GetString("server.router.home.name", "http.home"),
				Methods: services.GET,
				MiddleWares: []string{},
				Group: config.GetString("server.router.home.group", ""),
				Urls: []string{config.GetString("server.router.home.url", "/")},
				Handle: func(request services.Request, response services.Response) error {
					return response.View(config.GetString("server.router.home.name", "http.home"), map[string]interface{} {
					})
				},
			},
		})
	}
	return nil
}

func New(app services.Application, config services.Config, logger services.Logger, redis services.RedisClient, engine services.RenderEngine) services.Router {
	s := new(routerImpl)

	s.router = routing.New()
	s.server = &fasthttp.Server {
		Handler: s.router.HandleRequest,
		Name:    config.GetAsString("server.name", ""),
	}
	s.groups = make(map[string][]*services.RouteGroup)
	s.middleWares = make(map[string]*services.MiddleWare)
	s.log = logger
	s.engine = engine

	s.port = config.GetAsString("server.port", "8080")
	s.domain = config.GetAsString("server.address", "127.0.0.1")
	s.cert = config.GetString("server.connection.tsl.cert", "")
	s.key = config.GetString("server.connection.tsl.key", "")

	if config.GetBoolean("server.router.static.enabled", true) {
		s.fs = &fasthttp.FS {
			Root: config.GetString("server.router.static.path", app.ResourcesPath("")),
			Compress: true,
			AcceptByteRange: true,
		}
		s.fsHandler = s.fs.NewRequestHandler()
		s.router.Get("/assets/*", func(context *routing.Context) error {
			s.fsHandler(context.RequestCtx)
			return nil
		}).Name(config.GetString("server.router.static.name", "http.assets"))
	}
	if config.GetBoolean("server.sessions.enable", true) {
		switch config.GetAsString("server.sessions.storage.driver", "redis") {
		case "redis":
			sessions.UseDatabase(dialects.NewRedisDriver(redis.Connection(config.GetString("server.sessions.storage.connection", "default")), logger))
		}
		s.autoSession = true
	} else {
		s.autoSession = false
	}
	s.router.NotFound(func(context *routing.Context) error {
		s.log.ErrorFields("Method not found or allowed", map[string]interface{} {
			"request": context.Request.String(),
			"method": string(context.Method()),
		})
		return s.responseFromContext(context).View(config.GetString("server.router.error.name", "http.error"), map[string]interface{} {
			"message": "Method not found or allowed",
			"request": s.requestFromContext(context),
		})
	})
	sessions.UpdateConfig(sessions.Config{
		Cookie: config.GetString("server.sessions.cookie", "bahman"),
		Expires: time.Duration(config.GetInt("server.sessions.expires", 0)),
	})

	return s
}
