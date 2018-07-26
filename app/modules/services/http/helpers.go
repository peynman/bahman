package http

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/fasthttp-routing"
)

func (r *routerImpl) handleRoute(route *services.Route) func(context *routing.Context) error {
	return func(context *routing.Context) error {
		defer func() {
			if err := recover(); err != nil {
				r.log.InfoFields("Route recovering from error", map[string]interface{} {
					"route": route.Name,
					"err": err,
				})
				if appErr, ok := err.(services.ApplicationError); ok {
					r.log.InfoFields("Route recovered, answering with app error page", map[string]interface{} {
						"route": route.Name,
						"err": err,
					})
					r.responseFromContext(context).Error(appErr)
					return
				}
			}
		}()
		r.log.InfoFields("Handling route", map[string]interface{} {
			"route": route.Name,
		})
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

func (r *routerImpl) handleGroups(wares []*services.RouteGroup) []routing.Handler {
	if len(wares) > 0 {
		var methods = make([]routing.Handler, len(wares))
		for index, ware := range wares {
			methods[index] = r.handleMiddleWare(ware.Handler)
		}
		return methods
	}

	return []routing.Handler{r.handleEmpty}
}
func (r *routerImpl) handleMiddleWares(wares []*services.MiddleWare) []routing.Handler {
	if len(wares) > 0 {
		var methods = make([]routing.Handler, len(wares))
		for index, ware := range wares {
			methods[index] = r.handleMiddleWare(ware.Handler)
		}
		return methods
	}

	return []routing.Handler{r.handleEmpty}
}
func (r *routerImpl) handleMiddleWare(callback services.RequestHandler) routing.Handler {
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
	if method&services.GET != 0 {
		methods = append(methods, "GET")
	}
	if method&services.POST != 0 {
		methods = append(methods, "POST")
	}
	if method&services.PUT != 0 {
		methods = append(methods, "PUT")
	}
	if method&services.DELETE != 0 {
		methods = append(methods, "DELETE")
	}
	return
}