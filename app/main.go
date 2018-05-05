package main

import (
	"avalanche/app/libs/logger"
	"github.com/sirupsen/logrus"
	"avalanche/app/libs"
	"AvConfig/lib"
	"avalanche/app/libs/database"
	myConf "avalanche/app/libs/config"
)

func main() {
	config.Initialize(app.ConfigPath(""), app.RootPath(""), []config.EvaluatorFunction {
		new (myConf.StoragePathEvaluator),
		new (myConf.ResourcesPathEvaluator),
		new (myConf.RootPathEvaluator),
	})
	logger.Initialize()
	database.Initialize()
	defer database.Close()

	logger.InfoFields("Avalanche Server Start", logrus.Fields{
		"version": app.Version,
		"code": app.Code,
		"platform": app.Platform,
		"variant": app.Variant,
		"time": app.BuildTime,
	})

	//router := routing.New()
	//
	//router.Get("/", func(c *routing.Context) error {
	//	fmt.Fprintf(c, "Hello world")
	//	return nil
	//})
	//
	//panic(fasthttp.ListenAndServe(":"+config.GetAsString("server.port", ""), router.HandleRequest))
}


