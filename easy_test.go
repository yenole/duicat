package easy_go

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

type ReqLogin struct {
	Name string `json:"name"`
}

func (r ReqLogin) handleLogin(login *ReqLogin, token string, a int16) (interface{}, error) {
	panic(fmt.Errorf("error"))
	return 10000, nil
}

func (r ReqLogin) handleLogin2(token string, a int32) (interface{}, error) {
	return 111111, nil
}

func TestNewEasyGo(t *testing.T) {
	NewEasyGo().Get(`/api/`, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello")
	})

	NewEasyGo().Post(`/api/login`, JSON(ParseJSON(ReqLogin{}.handleLogin, (*ReqLogin)(nil), "token", "a")))

	NewEasyGo().Get(`/api/login2`, JSON(ParseParam(ReqLogin{}.handleLogin2, "token", "a")))

	http.ListenAndServe(":8080", nil)
}
