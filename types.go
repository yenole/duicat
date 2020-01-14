package easy_go

import "net/http"

type HandlerRenderFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)
