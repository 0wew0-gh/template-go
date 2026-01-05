package tool

import (
	"log"
	"net/http"
)

// 无防跨域(写入日志)
func (s Setting) PublicHandleNoAllowLog(w http.ResponseWriter, req *http.Request, log *log.Logger) {
	log.Println(">", req.Method, ":", req.Header.Get("X-Forwarded-For"), req.RemoteAddr, "->", req.RequestURI)
	if s.Config.Development {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
	}
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	w.Header().Set("content-type", "multipart/form-data")
	w.Header().Set("content-type", "application/json")
}
