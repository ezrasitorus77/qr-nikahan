package middleware

import (
	"net/http"
	"qr-nikahan/internal/helper"
)

type LogMiddleware struct {
	Handler http.Handler
}

func (m *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	helper.INFO(r.URL.Path)
	m.Handler.ServeHTTP(w, r)
	helper.INFO("DONE")
}
