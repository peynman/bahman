package cache

import "github.com/uniplaces/carbon"



func Put(key string, val interface{}, duration carbon.Carbon) bool {
	return true
}
func Get(key string, def interface{}) interface{} {
	return new(interface{})
}
func Exists(key string) bool {
	return true
}
func Forget(key string) {

}





