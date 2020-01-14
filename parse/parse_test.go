package parse

import (
	"fmt"
	"github.com/yenole/easy-go"
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

func TestParse(t *testing.T) {
	easy_go.NewEasyGo().Post(`/api/login`, easy_go.JSON(JSON(ReqLogin{}.handleLogin, (*ReqLogin)(nil), "token", "a")))

	easy_go.NewEasyGo().Get(`/api/login2`, easy_go.JSON(Param(ReqLogin{}.handleLogin2, "token", "a")))

	http.ListenAndServe(":8080", nil)
}
