package tool

import (
	"log"
	"net/http"
)

// 无防跨域(写入日志)
func PublicHandleNoAllowLog(w http.ResponseWriter, req *http.Request, log *log.Logger) {
	log.Println(">", req.Method, ":", req.Header.Get("X-Forwarded-For"), req.RemoteAddr, "->", req.RequestURI)
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	w.Header().Set("content-type", "multipart/form-data")
	w.Header().Set("content-type", "application/json")
}
