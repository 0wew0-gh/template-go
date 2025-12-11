package tool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kagurazakayashi/libNyaruko_Go/nyahttphandle"
	"github.com/kagurazakayashi/libNyaruko_Go/nyaredis"
)

var (
	mySQLCheckInterval time.Duration = 5 * time.Second // mysql检查间隔
)

// 测试连接
func (c Setting) MonitorMySQLConnection() {
	ticker := time.NewTicker(mySQLCheckInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := c.SQL.Ping()
		if err != nil {
			c.L.Printf(Error, "Ping failed: %v", err)
			os.Exit(0)
			break
		}
	}
}

// 错误返回
func (c Setting) BackErrorMsg(w http.ResponseWriter, localeID int, respID interface{}, errCode int, err interface{}, showErr bool, data interface{}) []byte {
	// 定义正则表达式
	re := regexp.MustCompile(`key '(.+)_UNIQUE'`)
	switch errType := err.(type) {
	case error:
		errStr := errType.Error()
		// fmt.Println(">>", errStr)
		match := re.FindStringSubmatch(errStr)
		// for _, vv := range match {
		// fmt.Println(">>", vv)
		// }
		if len(match) < 2 {
			break
		}
		errCode = 9008
		errStrList := strings.Split(match[1], ".")
		// fmt.Println(errStrList)
		if len(errStrList) < 2 {
			break
		}
		err = errStrList[1]
	case []string:
		errStr := ""
		for i, v := range errType {
			match := re.FindStringSubmatch(v)
			// for _, vv := range match {
			// fmt.Println(">>", vv)
			// }

			if len(match) < 2 {
				continue
			}
			errCode = 9008
			errStrList := strings.Split(match[1], ".")
			if len(errStrList) < 2 {
				continue
			}
			if i == 0 {
				errStr = ""
			}
			if errStr != "" {
				errStr += ","
			}
			errStr += errStrList[1]
		}
		if errStr != "" {
			err = errStr
		}
	}
	reData := map[string]interface{}{}
	reData["data"] = data
	c.L.Printf(Error, "[%v-%d]%+v", respID, errCode, err)
	if showErr {
		reData["err"] = err
		return nyahttphandle.AlertInfoJsonKV(w, localeID, respID, errCode, "", reData)
	}
	if data != nil {
		return nyahttphandle.AlertInfoJsonKV(w, localeID, respID, errCode, "", reData)
	}
	return nyahttphandle.AlertInfoJson(w, localeID, respID, errCode)
}

// 对比验证码
func (c Setting) VerifyVCode(user string, vcode string, isShowPrint bool) (bool, int, error) {
	temp := c.R_vcode.GetString(user, nyaredis.Option_isDelete(true))
	err := c.R_vcode.Error()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, -1, nil
		}
		return false, 9011, err
	}
	if temp != vcode {
		return false, -1, nil
	}
	return true, -1, nil
}

// 验证token
func (c Setting) VerifyToken(t string, isShowPrint bool) (map[string]interface{}, int, error) {
	userInfoStr := c.R_login.GetString("temp_" + t)
	err := c.R_login.Error()
	if err != nil {
		if strings.Contains(err.Error(), "redis: nil") {
			return nil, 3900, err
		}
		return nil, 9011, err
	}

	var userInfo map[string]interface{}
	err = json.Unmarshal([]byte(userInfoStr), &userInfo)
	if err != nil {
		return nil, 3901, err
	}

	return userInfo, -1, nil
}

// MARK: 组合字符串
func AddStr(s string, str string, splic string) string {
	if len(str) == 0 {
		return s
	}
	if s != "" {
		s += splic
	}
	s += str
	return s
}

// ==========
//
//	组装where语句
//	where		string	"原where语句"
//	new		string	"需要加入where的语句"
//	delimiter	string	"连接符"
//	isApostrophe	bool	"是否需要加入单引号"
func AssembleWhere(where string, new string, delimiter string, isApostrophe bool) string {
	if where != "" {
		where += delimiter
	}
	if isApostrophe {
		where += "'"
	}
	where += new
	if isApostrophe {
		where += "'"
	}
	return where
}

func GenerateOrderBy(fo []string, isho bool, fosc []string, ishosc bool, flt []string, ishlt bool, defaultby string, orderKey []string, orderBy []string, andorder string) (string, string, int, error) {
	orderSC := "ASC"
	order := "`" + defaultby + "` "
	limit := ""

	if ishosc {
		switch fosc[0] {
		case "desc", "DESC", "de", "De", "dE", "DE":
			orderSC = "DESC"
		}
	}
	if isho {
		for i := 0; i < len(orderKey); i++ {
			if fo[0] == orderKey[i] {
				order = "`" + orderBy[i] + "` "
			}
		}
	}
	order += orderSC
	if andorder != "" {
		if order != "" {
			order += ","
		}
		order += andorder
	}
	if ishlt {
		ls := strings.Split(flt[0], ",")
		if len(ls) < 1 || len(ls) > 2 {
			return "", "", 2041, fmt.Errorf("lt")
		}
		for _, v := range ls {
			_, err := strconv.Atoi(v)
			if err != nil {
				return "", "", 2041, fmt.Errorf("lt")
			}
		}
		switch len(ls) {
		case 1:
			limit = "0," + ls[0]
		case 2:
			limit = ls[0] + "," + ls[1]
		}
	}
	return order, limit, -1, nil
}
