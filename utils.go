package easy_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func JSON(handle HandlerRenderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stdout, stderr := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})

		defer func() {
			if stderr.Len() > 0 {
				_, _ = io.WriteString(w, fmt.Sprintf(`{"ret":false,"data":"%v"}`, stderr.String()))
			} else if stdout.Len() > 0 {
				_, _ = io.WriteString(w, fmt.Sprintf(`{"ret":true,"data":%v}`, stdout.String()))
			} else {
				_, _ = io.WriteString(w, `{"ret":true}`)
			}
			stdout.Reset()
			stderr.Reset()
		}()

		defer func() {
			if err := recover(); err != nil {
				stderr.WriteString(fmt.Sprint(err))
			}
		}()

		data, err := handle(w, r)
		if err != nil {
			stderr.WriteString(err.Error())
			return
		}

		if data == nil {
			return
		}

		byts, err := json.Marshal(data)
		if err != nil {
			stderr.WriteString(err.Error())
		} else {
			stdout.Write(byts)
		}

	}
}
