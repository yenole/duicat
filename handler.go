package duicat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

func RJSON(handle HandlerRenderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stdout, stderr := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})

		defer func() {
			w.Header().Add("content-type", "application/json")
			if stderr.Len() > 0 {
				_, _ = io.WriteString(w, fmt.Sprintf(`{"ret":false,"data":"%v"}`, stderr.String()))
			} else if stdout.Len() > 0 {
				_, _ = io.WriteString(w, fmt.Sprintf(`{"ret":true,"data":%v}`, stdout.String()))
			} else {
				_, _ = io.WriteString(w, `{"ret":true}`)
			}
			stdout.Reset()
			stderr.Reset()
		}()

		defer func() {
			if err := recover(); err != nil {
				stderr.WriteString(fmt.Sprint(err))
			}
		}()

		data, err := handle(w, r)
		if err != nil {
			stderr.WriteString(err.Error())
			return
		}

		if data == nil {
			return
		}

		byts, err := json.Marshal(data)
		if err != nil {
			stderr.WriteString(err.Error())
		} else {
			stdout.Write(byts)
		}

	}
}

func RPlain(handle HandlerRenderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer := bytes.NewBuffer([]byte{})

		defer func() {
			w.Header().Add("content-type", "text/plain")
			if buffer.Len() > 0 {
				_, _ = w.Write(buffer.Bytes())
			}
			buffer.Reset()
		}()

		defer func() {
			if err := recover(); err != nil {
				buffer.WriteString(fmt.Sprint(err))
			}
		}()

		data, err := handle(w, r)
		if err != nil {
			buffer.WriteString(err.Error())
			return
		}
		buffer.WriteString(fmt.Sprint(data))
	}
}

func P(handle interface{}, params ...string) HandlerRenderFunc {
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
		values, err := fetchParams(r, handleType, 0, params)
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

func PJSON(handle interface{}, json interface{}, params ...string) HandlerRenderFunc {
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

		values, err := fetchParams(r, handleType, 1, params)
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

func fetchParams(r *http.Request, el reflect.Type, skip int, params []string) ([]reflect.Value, error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	values := make([]reflect.Value, el.NumIn()-skip)
	for i := skip; i < el.NumIn(); i++ {
		if el.In(i).Kind() == reflect.String {
			values[i-skip] = reflect.ValueOf(query.Get(params[i-skip]))
		} else if len(query.Get(params[i-skip])) > 0 {
			var ret = reflect.New(el.In(i))
			err := json.Unmarshal([]byte(query.Get(params[i-skip])), ret.Interface())
			if err != nil {
				return nil, err
			}
			values[i-skip] = ret.Elem()
		} else {
			values[i-skip] = reflect.New(el.In(i)).Elem()
		}
	}
	return values, nil
}
