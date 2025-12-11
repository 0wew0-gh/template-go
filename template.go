package main

import (
	"net/http"
	"sync"
	"template-go/tool"

	"github.com/kagurazakayashi/libNyaruko_Go/nyahttphandle"
)

func template(w http.ResponseWriter, r *http.Request, c chan []byte) {
	tool.PublicHandleNoAllowLog(w, r, conf.L.Info)

	localeID := conf.Config.DefaultLocaleID

	if r.Method == http.MethodOptions {
		c <- nyahttphandle.AlertInfoJson(w, localeID, 99900, 1001)
		return
	} else if r.Method != http.MethodPost { // 检查是否为post请求
		// 返回 不是POST请求 的错误
		c <- nyahttphandle.AlertInfoJson(w, localeID, 99901, 2001)
		return
	}

	r.ParseMultipartForm(32 << 20)
	ft, isht := r.Form["t"] //token

	fshowErr, ishshowErr := r.Form["se"] //是否显示错误信息
	fl, ishl := r.Form["l"]              //语言
	localeID = conf.SetLanguage(ishl, fl)
	showErr := false
	if ishshowErr && fshowErr[0] == "1" {
		showErr = true
	}
	if !isht {
		c <- nyahttphandle.AlertInfoJsonKV(w, localeID, 11512, 2040, "p", "t")
		return
	}

	userInfo, errCode, err := conf.VerifyToken(ft[0], showErr)
	if err != nil {
		c <- conf.BackErrorMsg(w, localeID, 11513, errCode, err, showErr, nil)
		return
	}

	c <- nyahttphandle.AlertInfoJsonKV(w, localeID, 99999, 10000, "", userInfo)
}

func templateHandleFunc(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	c := make(chan []byte)
	go template(w, r, c)
	re := <-c
	wg.Done()
	w.Write([]byte(re))
	wg.Wait()
}
