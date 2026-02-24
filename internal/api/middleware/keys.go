package middleware

import (
	"net"
	"net/http"
)

func IPKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func UserKey(r *http.Request) string {
	user, ok := r.Context().Value(userContextKey).(userData)
	if !ok {
		return ""
	}
	return user.ID.String()
}
