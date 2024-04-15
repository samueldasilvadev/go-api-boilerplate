package apiGwAdapter

import "net/http"

type ResponseWriter struct {
	Status int
	Body   []byte
}

func (w *ResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return len(data), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
}
