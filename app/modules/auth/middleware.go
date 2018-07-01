package auth

import "github.com/peyman-abdi/avalanche/app/interfaces/services"

func (*AuthenticationModule) MiddleWares() []*services.MiddleWare {
	return []*services.MiddleWare {
		{
			Name: "Auth",
			Priority: 1,
			Handler: func(request services.Request, response services.Response) error {

				return nil
			},
		},
		{
			Name: "Guest",
			Priority: 1,
			Handler: func(request services.Request, response services.Response) error {
				return nil
			},
		},
	}
}