package main

import (
	"fmt"
	"github.com/yenole/duicat"
	"github.com/yenole/duicat/parse"
	"github.com/yenole/duicat/render"
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
		Get("h2", render.JSON(handler.h2)).
		Get("h3", render.JSON(parse.Param(handler.h3, "name", "age"))).
		Post("json", render.JSON(parse.JSON(handler.json, (*login)(nil), "token"))).
		Get("err", render.JSON(parse.Param(handler.err, "token")))

	dc.HandleGroup("/v2", duicat.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("v2"))
	})

	_ = dc.Listen(":8081")
}
