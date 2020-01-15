package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yenole/duicat"
	"io"
	"net/http"
)

func JSON(handle duicat.HandlerRenderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stdout, stderr := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})

		defer func() {
			w.Header().Add("content-type", "application/json")
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

func Plain(handle duicat.HandlerRenderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer := bytes.NewBuffer([]byte{})

		defer func() {
			w.Header().Add("content-type", "text/plain")
			if buffer.Len() > 0 {
				_, _ = w.Write(buffer.Bytes())
			}
			buffer.Reset()
		}()

		defer func() {
			if err := recover(); err != nil {
				buffer.WriteString(fmt.Sprint(err))
			}
		}()

		data, err := handle(w, r)
		if err != nil {
			buffer.WriteString(err.Error())
			return
		}
		buffer.WriteString(fmt.Sprint(data))
	}
}
