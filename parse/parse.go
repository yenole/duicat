package parse

import (
	"encoding/json"
	"github.com/yenole/easy-go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

func Param(handle interface{}, params ...string) easy_go.HandlerRenderFunc {
	handleType := reflect.TypeOf(handle)
	if handleType.Kind() != reflect.Func {
		panic("handle is not a function")
	}

	if handleType.NumIn() != len(params) {
		panic("number of handles and params do not match")
	}

	if handleType.NumOut() != 2 || handleType.Out(0).Kind() != reflect.Interface && handleType.Out(1).Kind() != reflect.Interface {
		panic("handle return parameter error")
	}

	handleValue := reflect.ValueOf(handle)

	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		values, err := fetchParams(r, handleType, 0, params...)
		if err != nil {
			return nil, err
		}

		results := handleValue.Call(values)

		if results[1].IsNil() {
			return results[0].Interface(), nil
		}

		return nil, results[1].Interface().(error)
	}
}

func JSON(handle interface{}, json interface{}, params ...string) easy_go.HandlerRenderFunc {
	handleType := reflect.TypeOf(handle)
	if handleType.Kind() != reflect.Func {
		panic("handle is not a function")
	}

	if handleType.NumIn() != len(params)+1 {
		panic("number of handles and params do not match")
	}

	if handleType.NumOut() != 2 || handleType.Out(0).Kind() != reflect.Interface && handleType.Out(1).Kind() != reflect.Interface {
		panic("handle return parameter error")
	}

	jsonType := reflect.TypeOf(json)
	handleValue := reflect.ValueOf(handle)

	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		data := reflect.New(jsonType)

		err := fetchBody2Json(r.Body, data.Interface())
		if err != nil {
			return nil, err
		}

		values, err := fetchParams(r, handleType, 1, params...)
		if err != nil {
			return nil, err
		}

		results := handleValue.Call(append([]reflect.Value{data.Elem()}, values...))

		if results[1].IsNil() {
			return results[0].Interface(), nil
		}

		return nil, results[1].Interface().(error)
	}
}

func fetchBody2Json(body io.ReadCloser, raw interface{}) error {
	byts, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(byts, raw)
}

func fetchParams(r *http.Request, el reflect.Type, skip int, params ...string) ([]reflect.Value, error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	values := make([]reflect.Value, el.NumIn()-skip)
	for i := skip; i < el.NumIn(); i++ {
		if el.In(i).Kind() == reflect.String {
			values[i-skip] = reflect.ValueOf(query.Get(params[i-skip]))
		} else {
			var ret = reflect.New(el.In(i))
			err := json.Unmarshal([]byte(query.Get(params[i-skip])), ret.Interface())
			if err != nil {
				return nil, err
			}
			values[i-skip] = ret.Elem()
		}

	}
	return values, nil
}
