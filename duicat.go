package duicat

import "net/http"

type DuiCat struct {
	*http.ServeMux
}

func NewDuiCat() *DuiCat {
	return &DuiCat{ServeMux: http.DefaultServeMux}
}

func (easy *DuiCat) AddGroup(path string) *DuiCat {
	return easy
}

func (easy *DuiCat) Post(path string, handle func(w http.ResponseWriter, r *http.Request)) *DuiCat {
	easy.HandleFunc(path, handle)
	return easy
}

func (easy *DuiCat) Get(path string, handle func(w http.ResponseWriter, r *http.Request)) *DuiCat {
	easy.HandleFunc(path, handle)
	return easy
}
