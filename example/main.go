package main

import (
	"fmt"
	"github.com/yenole/duicat"
	"net/http"
)

type login struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type handler struct {
}

func (handler) h1(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("h1"))
}

func (handler) h2(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return r.URL.Path, nil
}

func (handler) h3(name string, age uint) (interface{}, error) {
	return fmt.Sprint("h3.", name, ":", age), nil
}

func (handler) json(login *login, token string) (interface{}, error) {
	fmt.Println(login, token)
	return login, nil
}

func (handler) err(token string) (interface{}, error) {
	if len(token) == 0 {
		panic(fmt.Errorf("token is empty"))
	}
	return nil, fmt.Errorf(token)
}

func main() {
	handler := handler{}

	dc := duicat.NewDuiCat()
	dc.Group("/v1").
		Get("h1", handler.h1).
		Get("h2", duicat.RJSON(handler.h2)).
		Get("h3", duicat.RJSON(duicat.P(handler.h3, "name", "age"))).
		Post("json", duicat.RJSON(duicat.PJSON(handler.json, (*login)(nil), "token"))).
		Get("err", duicat.RJSON(duicat.PJSON(handler.err, "token")))

	dc.HandleGroup("/v2", duicat.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("v2"))
	})

	_ = dc.Listen(":8081")
}
