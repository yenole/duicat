package duicat

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type RenderHandler func(w http.ResponseWriter, r *http.Request) interface{}

func handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, nil
}

func Println(a interface{}) {
	_, ok := a.(*RenderHandler)
	value := reflect.ValueOf(a).Convert(reflect.TypeOf(RenderHandler(nil)))
	newhandle := value.Interface().(RenderHandler)
	fmt.Println(a, ok, value, newhandle)
}

func TestTest(t *testing.T) {
	Println(handle)
}
