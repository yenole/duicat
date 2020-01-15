package duicat

import (
	"fmt"
	"net/http"
	"strings"
)

type DuiCat struct {
	handler Handler
}

func NewDuiCat() *DuiCat {
	return &DuiCat{handler: http.DefaultServeMux}
}

func (dc *DuiCat) Listen(addr string) error {
	return http.ListenAndServe(addr, dc.handler)
}

func (dc *DuiCat) HandleFunc(path string, method Method, handle http.HandlerFunc) *DuiCat {
	dc.handler.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		if method == MethodAll || request.Method == string(method) {
			handle(writer, request)
		}
	})
	return dc
}

func (dc *DuiCat) Post(pattern string, handle func(w http.ResponseWriter, r *http.Request)) *DuiCat {
	return dc.HandleFunc(pattern, MethodPost, handle)
}

func (dc *DuiCat) Get(pattern string, handle func(w http.ResponseWriter, r *http.Request)) *DuiCat {
	return dc.HandleFunc(pattern, http.MethodGet, handle)
}

func (dc *DuiCat) Group(pattern string) *Group {
	if !strings.HasSuffix(pattern, "/") {
		pattern = fmt.Sprint(pattern, "/")
	}
	return &Group{prefix: pattern, handler: dc}
}

func (dc *DuiCat) HandleGroup(pattern string, method Method, handle http.HandlerFunc) *Group {
	return dc.HandleFunc(pattern, method, handle).Group(pattern)
}

type Group struct {
	prefix  string
	handler *DuiCat
}

func (gp *Group) HandleFunc(path string, method Method, handle http.HandlerFunc) *Group {
	gp.handler.HandleFunc(fmt.Sprint(gp.prefix, path), method, handle)
	return gp
}

func (gp *Group) Post(pattern string, handle func(w http.ResponseWriter, r *http.Request)) *Group {
	gp.handler.Post(fmt.Sprint(gp.prefix, pattern), handle)
	return gp
}

func (gp *Group) Get(pattern string, handle func(w http.ResponseWriter, r *http.Request)) *Group {
	gp.handler.Get(fmt.Sprint(gp.prefix, pattern), handle)
	return gp
}
