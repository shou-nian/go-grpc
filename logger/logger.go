package logger

import (
	"io"
	"log/slog"
	"net/http"
)

func PrintRequestInfo(req *http.Request) {
	method := req.Method
	url := req.URL
	var b []byte
	body := req.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			slog.Error("close request body failed: ", err)
			return
		}
	}(body)

	_, err := body.Read(b)
	if err != nil {
		slog.Info("request info: ", method, " ", url)
	}

	slog.Info("request info: ", method, " ", url, " ", b)
}

func ErrorLogger(err error) {
	slog.Info(err.Error())
}
