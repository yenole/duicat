package easy_go

import "net/http"

type EasyGo struct {
	*http.ServeMux
}

func NewEasyGo() *EasyGo {
	return &EasyGo{ServeMux: http.DefaultServeMux}
}

func (easy *EasyGo) AddGroup(path string) *EasyGo {
	return easy
}

func (easy *EasyGo) Post(path string, handle func(w http.ResponseWriter, r *http.Request)) *EasyGo {
	easy.HandleFunc(path, handle)
	return easy
}

func (easy *EasyGo) Get(path string, handle func(w http.ResponseWriter, r *http.Request)) *EasyGo {
	easy.HandleFunc(path, handle)
	return easy
}
