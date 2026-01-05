package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"template-go/tool"

	"github.com/kagurazakayashi/libNyaruko_Go/nyamysql"
	"github.com/kagurazakayashi/libNyaruko_Go/nyaredis"
)

// / 语言
// / 1: en
// / 2: zhHans
// / 3: zhHant
// / 4: es
// / ====================
var (
	version string = "0.0.1"

	conf tool.Setting // 设置

	updateMailLogFile *os.File = nil
)

func main() {
	var (
		tn  time.Time = time.Now()
		err error     = nil
	)

	//获取设置
	conf, err = tool.GetPublicVariable("./conf/setting")
	if err != nil {
		fmt.Println("读取配置文件失败:", err)
		return
	}

	tool.SetupCloseHandler()

	conf.L, updateMailLogFile = conf.L.SetupLogFile(conf.Config.LogPath, updateMailLogFile)

	// 输出版本信息
	conf.L.Printf(tool.Info, "[%s] %s: v%s\n", tn.In(conf.Config.TimeZone.Location).Format(conf.Config.TimeZone.Format), conf.Config.AppName, version)

	if conf.Config.LoggerLevel < tool.Debug {
		conf.L.Println(tool.Info, "MySQL no debug log")
		conf.SQL = nyamysql.NewC(conf.DBSetting.Mysql, nil, conf.DBSetting.LoggerLevel)
	} else {
		conf.L.Println(tool.Debug, "MySQL debug test")
		conf.SQL = nyamysql.NewC(conf.DBSetting.Mysql, conf.L.Debug, conf.DBSetting.LoggerLevel)
	}
	if conf.SQL.Error() != nil {
		conf.L.Println(tool.Error, "MySQL Link failed:", conf.SQL.Error())
		return
	}
	conf.SQL.SetMaxOpenConns(conf.DBSetting.MysqlAdvanced.MaxLinkNumber)
	conf.SQL.SetMaxIdleConns(conf.DBSetting.MysqlAdvanced.MaxIdleNumber)
	conf.SQL.SetConnMaxIdleTime(time.Duration(conf.DBSetting.MysqlAdvanced.MaxIdleTime) * time.Minute)
	conf.SQL.SetConnMaxLifetime(time.Duration(conf.DBSetting.MysqlAdvanced.MaxLifeTime) * time.Minute)

	go conf.MonitorMySQLConnection()

	err = conf.GetOssConfig()
	if err != nil {
		conf.L.Println(tool.Error, "GetOssConfig err:", err)
		return
	}

	conf.L.Println(tool.Info, "Redis Link test")
	conf.R_vcode = nyaredis.NewC(conf.DBSetting.Redis, conf.Config.RedisVCodeDB, conf.DBSetting.RedisMaxDB)
	if conf.R_vcode.Error() != nil {
		conf.L.Println(tool.Error, "Redis Vcode Link failed:", conf.R_vcode.Error())
		return
	}
	conf.R_login = nyaredis.NewC(conf.DBSetting.Redis, conf.Config.RedisLoginDB, conf.DBSetting.RedisMaxDB)
	if conf.R_vcode.Error() != nil {
		conf.L.Println(tool.Error, "Redis Login Link failed:", conf.R_login.Error())
		return
	}
	conf.L.Println(tool.Info, "Redis test link success")

	// 定时处理
	go AutomaticHandle(tn)

	conf.L.Println(tool.Info, "监听端口:", conf.Config.ListenAndServe)
	conf.L.Println(tool.Info, "初始化完成")
	conf.L.Println(tool.Info, "====================")

	var handleFuncList []func(http.ResponseWriter, *http.Request)

	// test 999
	handleFuncList = append(handleFuncList, templateHandleFunc)

	if len(handleFuncList) != len(conf.Config.HandleFuncKeyList) {
		conf.L.Println(tool.Info, len(handleFuncList), len(conf.Config.HandleFuncKeyList))
		conf.L.Println(tool.Info, "接口数量不匹配")
		return
	}

	http.HandleFunc(conf.Config.SubURL, mainHandleFunc)
	conf.L.Println(tool.Info, "接口列表:")
	conf.L.Println(tool.Info, conf.Config.SubURL)
	for i := 0; i < len(handleFuncList); i++ {
		pattern := conf.Config.SubURL + conf.Config.HandleFuncKeyList[i]
		conf.L.Println(tool.Info, pattern)
		http.HandleFunc(pattern, handleFuncList[i])
	}

	err = http.ListenAndServe(":"+conf.Config.ListenAndServe, nil)
	if err != nil {
		conf.L.Println(tool.Error, "ListenAndServe error:", err)
		return
	}
}

func mainHandleFunc(w http.ResponseWriter, r *http.Request) {
	conf.PublicHandleNoAllowLog(w, r, conf.L.Info)
	w.WriteHeader(404)
	w.Write([]byte{})
}
