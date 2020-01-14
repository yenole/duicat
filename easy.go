package easy_go

import "net/http"

type EasyGo struct {
}

func NewEasyGo() *EasyGo {
	return &EasyGo{}
}

func (easy *EasyGo) AddGroup(path string) *EasyGo {
	return easy
}

func (easy *EasyGo) Post(path string, handle func(w http.ResponseWriter, r *http.Request)) *EasyGo {
	return easy
}

func (easy *EasyGo) Get(path string, handle func(w http.ResponseWriter, r *http.Request)) *EasyGo {
	return easy
}
