package easy_go

import "net/http"

func ParamJSON(handle func(data interface{}) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (i interface{}, err error) {
		return nil, nil
	}
}

func JSON(handle func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func Template(tmpl string, handle interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func Plain() {

}

func HTML() {

}
