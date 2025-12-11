package main

import (
	"template-go/tool"
	"time"
)

func AutomaticHandle(tn time.Time) {
	// 对准时间
	next := tn.Truncate(time.Minute).Add(time.Minute)
	time.Sleep(time.Until(next))

	// 创建Log定时器
	var tickerLog *time.Ticker = time.NewTicker(time.Minute * time.Duration(1))

	// 创建定时器
	var ticker *time.Ticker = time.NewTicker(6 * time.Hour)
	for {
		select {
		case <-tickerLog.C:
			conf.L, updateMailLogFile = conf.L.SetupLogFile(conf.Config.LogPath, updateMailLogFile)
		case <-ticker.C:
			time.Sleep(6 * time.Hour)
			conf.L.Println(tool.Info, "====================")
			conf.L.Println(tool.Info, "AutomaticHandleGuest")
		}
	}
}
