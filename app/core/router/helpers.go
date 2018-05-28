package router

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/peyman-abdi/avalanche/app/interfaces"
)

func handleNotFound(c *routing.Context) error {
	c.Error("Page not found!", 404)
	return nil
}

func handleRoute(route *interfaces.Route) func(context *routing.Context) error {
	return func(context *routing.Context) error {
		if route.Verify != nil {
			if err := route.Verify(requestFromContext(context)); err != nil {
				return err
			}
		}

		if route.Handle != nil {
			return route.Handle(requestFromContext(context), responseFromContext(context))
		}

		return nil
	}
}

func handleMiddleWares(callback []interfaces.RequestHandler) []routing.Handler {
	if len(callback) > 0 {
		var methods = make([]routing.Handler, len(callback))
		for index, method := range callback {
			methods[index] = handleMiddleWare(method)
		}
		return methods
	}

	return []routing.Handler {handleEmpty}
}
func handleMiddleWare(callback interfaces.RequestHandler) routing.Handler {
	if callback != nil {
		return func(context *routing.Context) error {
			return callback(requestFromContext(context), responseFromContext(context))
		}
	}

	return handleEmpty
}

func handleEmpty(context *routing.Context) error {
	return nil
}

func methodsFromInt(method int) (methods []string) {
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

func requestFromContext(ctx *routing.Context) interfaces.Request {
	return nil
}
func responseFromContext(ctx *routing.Context) interfaces.Response {
	return nil
}