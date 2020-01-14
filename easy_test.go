package easy_go

import (
	"io"
	"net/http"
	"testing"
)

func TestNewEasyGo(t *testing.T) {
	NewEasyGo().Post(`/api/`, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello")
	})

	NewEasyGo().Post(`/api/login`, JSON(ParamJSON(func(data interface{}) (i interface{}, err error) {
		return nil, nil
	})))

	NewEasyGo().Post(`/api/login2`, JSON(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return nil, nil
	}))
}
