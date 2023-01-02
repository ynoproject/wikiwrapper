package common

import "net/http"

func GetIp(r *http.Request) string {
	return r.Header.Get("x-forwarded-for")
}
