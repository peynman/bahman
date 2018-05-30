package router

import (
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

func handleNotFound(c *routing.Context) error {
	c.Error("Page not found!", 404)
	return nil
}

func (r *routerImpl) handleRoute(route *interfaces.Route) func(context *routing.Context) error {
	return func(context *routing.Context) error {
		if route.Verify != nil {
			if err := route.Verify(r.requestFromContext(context)); err != nil {
				return err
			}
		}

		if route.Handle != nil {
			return route.Handle(r.requestFromContext(context), r.responseFromContext(context))
		}

		return nil
	}
}

func (r *routerImpl) handleGroups(wares []*interfaces.RouteGroup) []routing.Handler {
	if len(wares) > 0 {
		var methods = make([]routing.Handler, len(wares))
		for index, ware := range wares {
			methods[index] = r.handleMiddleWare(ware.Handler)
		}
		return methods
	}

	return []routing.Handler {r.handleEmpty}
}
func (r *routerImpl) handleMiddleWares(wares []*interfaces.MiddleWare) []routing.Handler {
	if len(wares) > 0 {
		var methods = make([]routing.Handler, len(wares))
		for index, ware := range wares {
			methods[index] = r.handleMiddleWare(ware.Handler)
		}
		return methods
	}

	return []routing.Handler {r.handleEmpty}
}
func (r *routerImpl) handleMiddleWare(callback interfaces.RequestHandler) routing.Handler {
	if callback != nil {
		return func(context *routing.Context) error {
			return callback(r.requestFromContext(context), r.responseFromContext(context))
		}
	}

	return r.handleEmpty
}

func (r *routerImpl) handleEmpty(context *routing.Context) error {
	return nil
}

func (r *routerImpl) methodsFromInt(method int) (methods []string) {
	if method & interfaces.GET != 0 {
		methods = append(methods, "GET")
	}
	if method & interfaces.POST != 0 {
		methods = append(methods, "POST")
	}
	if method & interfaces.PUT != 0 {
		methods = append(methods, "PUT")
	}
	if method & interfaces.DELETE != 0 {
		methods = append(methods, "DELETE")
	}
	return
}
func (r *routerImpl) responseFromContext(ctx *routing.Context) interfaces.Response {
	return &responseImpl{ context: ctx, log: r.log }
}