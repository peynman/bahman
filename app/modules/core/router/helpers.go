package router

import (
	"github.com/peyman-abdi/fasthttp-routing"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
)

func handleNotFound(c *routing.Context) error {
	c.Error("Page not found!", 404)
	return nil
}

func (r *routerImpl) handleRoute(route *core.Route) func(context *routing.Context) error {
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

func (r *routerImpl) handleGroups(wares []*core.RouteGroup) []routing.Handler {
	if len(wares) > 0 {
		var methods = make([]routing.Handler, len(wares))
		for index, ware := range wares {
			methods[index] = r.handleMiddleWare(ware.Handler)
		}
		return methods
	}

	return []routing.Handler{r.handleEmpty}
}
func (r *routerImpl) handleMiddleWares(wares []*core.MiddleWare) []routing.Handler {
	if len(wares) > 0 {
		var methods = make([]routing.Handler, len(wares))
		for index, ware := range wares {
			methods[index] = r.handleMiddleWare(ware.Handler)
		}
		return methods
	}

	return []routing.Handler{r.handleEmpty}
}
func (r *routerImpl) handleMiddleWare(callback core.RequestHandler) routing.Handler {
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
	if method&core.GET != 0 {
		methods = append(methods, "GET")
	}
	if method&core.POST != 0 {
		methods = append(methods, "POST")
	}
	if method&core.PUT != 0 {
		methods = append(methods, "PUT")
	}
	if method&core.DELETE != 0 {
		methods = append(methods, "DELETE")
	}
	return
}
func (r *routerImpl) responseFromContext(ctx *routing.Context) core.Response {
	return &responseImpl{context: ctx, log: r.log, engine: r.engine}
}
