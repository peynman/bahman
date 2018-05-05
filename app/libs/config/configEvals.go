package config

import (
	"avalanche/app/libs"
	"AvConfig/lib"
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
var _ config.EvaluatorFunction = (*StoragePathEvaluator)(nil)


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
var _ config.EvaluatorFunction = (*ResourcesPathEvaluator)(nil)


/**
	Root Path: resource(path)
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
var _ config.EvaluatorFunction = (*RootPathEvaluator)(nil)