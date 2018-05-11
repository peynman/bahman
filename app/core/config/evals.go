package config

import (
	gConfig "peyman/config"
	"avalanche/app/core/app"
	"time"
)

/**
	Storage Path: storage(path)
 */
type StoragePathEvaluator struct {
}
func (_ *StoragePathEvaluator) GetFunctionName() string {
	return "storage"
}
func (_ *StoragePathEvaluator) Eval(params []string, def interface{}) interface{} {
	if len(params) > 0 {
		return app.StoragePath(params[0])
	} else {
		return app.StoragePath("")
	}
}
var _ gConfig.EvaluatorFunction = (*StoragePathEvaluator)(nil)


/**
	Resource Path: resource(path)
 */
type ResourcesPathEvaluator struct {
}
func (_ *ResourcesPathEvaluator) GetFunctionName() string {
	return "resources"
}
func (_ *ResourcesPathEvaluator) Eval(params []string, def interface{}) interface{} {
	if len(params) > 0 {
		return app.ResourcesPath(params[0])
	} else {
		return app.ResourcesPath("")
	}
}
var _ gConfig.EvaluatorFunction = (*ResourcesPathEvaluator)(nil)


/**
	Root Path: root(path)
 */
type RootPathEvaluator struct {
}
func (_ *RootPathEvaluator) GetFunctionName() string {
	return "root"
}
func (_ *RootPathEvaluator) Eval(params []string, def interface{}) interface{} {
	if len(params) > 0 {
		return app.RootPath(params[0])
	} else {
		return app.RootPath("")
	}
}
var _ gConfig.EvaluatorFunction = (*RootPathEvaluator)(nil)


/**
	System parameters: system(parameter, default)
 */
type SystemParameterEvaluator struct {
}
func (_ *SystemParameterEvaluator) GetFunctionName() string {
	return "system"
}
func (_ *SystemParameterEvaluator) Eval(params []string, def interface{}) interface{} {
	if len(params) > 0 {
		switch params[0] {
		case "os":
		case "platform":
			return app.Platform
		case "variant":
			return app.Variant
		}
	}
	return def
}
var _ gConfig.EvaluatorFunction = (*SystemParameterEvaluator)(nil)


/**
	Time : time(parameter)
	parameter:
		hour, minute, second, now
 */
type TimeEvaluator struct {
}
func (_ *TimeEvaluator) GetFunctionName() string {
	return "time"
}
func (_ *TimeEvaluator) Eval(params []string, def interface{}) interface{} {
	if len(params) > 0 {
		switch params[0] {
		case "hour":
			return time.Hour
		case "minute":
			return time.Minute
		case "second":
			return time.Second
		case "now":
			return time.Now().String()
		}
	}
	return def
}
var _ gConfig.EvaluatorFunction = (*TimeEvaluator)(nil)