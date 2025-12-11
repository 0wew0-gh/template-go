package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kagurazakayashi/libNyaruko_Go/nyahttphandle"
	"github.com/kagurazakayashi/libNyaruko_Go/nyaio"
	"github.com/kagurazakayashi/libNyaruko_Go/nyamysql"
	"github.com/kagurazakayashi/libNyaruko_Go/nyaredis"
	"gopkg.in/yaml.v2"
)

type Setting struct {
	DBSetting DBSetting `json:"dbSetting" yaml:"dbSetting"`
	Config    Config    `json:"config" yaml:"config"`
	OssConfig *OssConfigSync
	SQL       *nyamysql.NyaMySQL
	R_vcode   *nyaredis.NyaRedis
	R_login   *nyaredis.NyaRedis
	L         Logger
}

type DBSetting struct {
	Mysql         nyamysql.MySQLDBConfig `json:"mysql" yaml:"mysql"`
	MysqlAdvanced MySQLAdvanced          `json:"mysql_advanced" yaml:"mysql_advanced"`
	Redis         nyaredis.RedisDBConfig `json:"redis" yaml:"redis"`
	RedisMaxDB    int                    `json:"redis_max_db" yaml:"redis_max_db"`
	MaxLink       MaxLinkNumber          `json:"maxLinkNumber" yaml:"maxLinkNumber"`
	WaitCount     int                    `json:"waitCount" yaml:"waitCount"`
	WaitTime      int                    `json:"waitTime" yaml:"waitTime"`
	LoggerLevel   int                    `json:"logger_level" yaml:"logger_level"`
}
type MySQLAdvanced struct {
	MaxLinkNumber int `json:"maxLinkNumber" yaml:"maxLinkNumber"`
	MaxIdleNumber int `json:"maxIdleNumber" yaml:"maxIdleNumber"`
	MaxIdleTime   int `json:"maxIdleTime" yaml:"maxIdleTime"`
	MaxLifeTime   int `json:"maxLifeTime" yaml:"maxLifeTime"`
	MaxPoolSize   int `json:"maxPoolSize" yaml:"maxPoolSize"`
	TimeOut       int `json:"timeOut" yaml:"timeOut"`
}
type MaxLinkNumber struct {
	Mysql int `json:"mysql"`
	Redis int `json:"redis"`
}

type Config struct {
	AppName           string   `json:"appName" yaml:"appName"`
	Development       bool     `json:"development" yaml:"development"`
	ReturnMsgFilePath string   `json:"returnMessageFilePath" yaml:"returnMessageFilePath"`
	TimeZone          TimeZone `json:"timeZone" yaml:"timeZone"`
	DefaultLocaleID   int      `json:"defaultLocaleID" yaml:"defaultLocaleID"` // 默认语言
	Language          []string `json:"lang" yaml:"lang"`                       // 语言列表
	LogPath           string   `json:"logPath" yaml:"logPath"`
	SubURL            string   `json:"suburl" yaml:"suburl"`                       // 子目录
	HandleFuncKeyList []string `json:"handleFuncKeyList" yaml:"handleFuncKeyList"` // 接口名称列表
	ListenAndServe    string   `json:"listenandserve" yaml:"listenandserve"`       // 监听端口
	RedisVCodeDB      int      `json:"redis_vcode_db" yaml:"redis_vcode_db"`       // 验证码数据库
	RedisLoginDB      int      `json:"redis_login_db" yaml:"redis_login_db"`       // 登录数据库
	TokenSaveTime     int      `json:"token_save_time" yaml:"token_save_time"`     // token保存时间
	///	日志等级
	///	-1: 禁用
	///	 0: Error
	///	 1: Warning
	///	 2: Debug
	///	 3: Infox
	LoggerLevel    int               `json:"logger_level" yaml:"logger_level"`
	Proxy          string            `json:"http_proxy" yaml:"http_proxy"` // http代理
	OSSDelPathList map[string]string `json:"oss_path" yaml:"oss_path"`     // OSS清理列表
	OSS            []OSSConfig       `json:"oss" yaml:"oss"`
	Mail           MailSetting       `json:"mail" yaml:"mail"`
}
type TimeZone struct {
	Location    *time.Location
	Format      string `json:"format" yaml:"format"`
	Offset      int    `json:"offset" yaml:"offset"`
	Name        string `json:"name" yaml:"name"`
	SQLLocation *time.Location
	SQLOffset   int    `json:"sqlOffset" yaml:"sqlOffset"`
	SQLName     string `json:"sqlName" yaml:"sqlName"`
} // 时区
type OSSConfig struct {
	Bucket          string `json:"bucket" yaml:"bucket"`
	AccessKeyID     string `json:"id" yaml:"id"`
	AccessKeySecret string `json:"secret" yaml:"secret"`
}
type MailSetting struct {
	From     string `json:"from" yaml:"from"`         //发送地址
	FormName string `json:"formname" yaml:"formname"` //发送地址名称
	PW       string `json:"pw" yaml:"pw"`             //密码
	Host     string `json:"host" yaml:"host"`         //发送邮箱host
	Port     string `json:"port" yaml:"port"`         //发送邮箱端口
	Logo     string `json:"logo" yaml:"logo"`         //邮件logo
	CSS      string `json:"css" yaml:"css"`           //邮件css
	Body     string `json:"body" yaml:"body"`         //邮件主体
	Type     string `json:"mailtype" yaml:"mailtype"` //邮件类型
}

func GetPublicVariable(path string) (Setting, error) {
	var (
		confFileType string = "YAML"
		conf         Setting
		confStr      string
		cstSh        *time.Location
		err          error
	)
	confStr, err = nyaio.FileRead(fmt.Sprintf("%s.yaml", path))
	if err != nil {
		confStr, err = nyaio.FileRead(fmt.Sprintf("%s.json", path))
		if err != nil {
			return conf, err
		}
		confFileType = "JSON"
	}
	switch confFileType {
	case "YAML":
		err = yaml.Unmarshal([]byte(confStr), &conf)
		if err != nil {
			return conf, err
		}
	case "JSON":
		err = json.Unmarshal([]byte(confStr), &conf)
		if err != nil {
			return conf, err
		}
	default:
		return conf, fmt.Errorf("配置文件类型错误")
	}

	if conf.DBSetting.Mysql.DbName == "" {
		return conf, fmt.Errorf("数据库配置错误")
	}
	if conf.DBSetting.MysqlAdvanced.MaxLinkNumber == 0 {
		conf.DBSetting.MysqlAdvanced.MaxLinkNumber = 10
		conf.DBSetting.MysqlAdvanced.MaxIdleNumber = 10
		conf.DBSetting.MysqlAdvanced.MaxIdleTime = 10
		conf.DBSetting.MysqlAdvanced.MaxLifeTime = 10
		conf.DBSetting.MysqlAdvanced.MaxPoolSize = 10
		conf.DBSetting.MysqlAdvanced.TimeOut = 120
	}
	if conf.DBSetting.WaitCount == 0 {
		conf.DBSetting.WaitCount = 10
	}
	if conf.DBSetting.WaitTime == 0 {
		conf.DBSetting.WaitTime = 500
	}
	if conf.DBSetting.LoggerLevel < nyamysql.NYAMYSQL_LOG_LEVEL_ERROR {
		conf.DBSetting.LoggerLevel = nyamysql.NYAMYSQL_LOG_LEVEL_ERROR
	}
	if conf.DBSetting.Redis.Address == "" {
		return conf, fmt.Errorf("Redis配置错误")
	}
	cstSh, err = time.LoadLocation(conf.Config.TimeZone.Name)
	if err != nil {
		// fmt.Println("时区文件加载失败:", err)
		cstSh = time.FixedZone("CST", conf.Config.TimeZone.Offset*3600)
		// fmt.Println("按时区加载")
		err = nil
	}
	conf.Config.TimeZone.Location = cstSh
	sql_cstSh, err := time.LoadLocation(conf.Config.TimeZone.SQLName)
	if err != nil {
		// fmt.Printf("[%s]时区文件加载失败:%v\n", conf.Config.TimeZone.SQLName, err)
		sql_cstSh = time.FixedZone("CST", conf.Config.TimeZone.SQLOffset*3600)
		// fmt.Println("按时区加载")
		err = nil
	}
	conf.Config.TimeZone.SQLLocation = sql_cstSh

	reint := nyahttphandle.AlertInfoTemplateLoad(conf.Config.ReturnMsgFilePath)
	if reint == 0 {
		return conf, fmt.Errorf("返回信息文件加载失败")
	}
	nyahttphandle.SetSuccessRange(10000, 10000)

	if conf.Config.LoggerLevel < -1 || conf.Config.LoggerLevel > 10 {
		// fmt.Println("Logger level error:", conf.Config.LoggerLevel)
		conf.Config.LoggerLevel = -1
	}
	conf.L.level = conf.Config.LoggerLevel

	_, ok := conf.Config.OSSDelPathList["home"]
	if !ok {
		conf.Config.OSSDelPathList["home"] = "home/"
	}
	_, ok = conf.Config.OSSDelPathList["news"]
	if !ok {
		conf.Config.OSSDelPathList["news"] = "news/"
	}
	_, ok = conf.Config.OSSDelPathList["product"]
	if !ok {
		conf.Config.OSSDelPathList["product"] = "product/"
	}

	conf.OssConfig = NewOssConfigSync(nil)
	return conf, err
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// fmt.Println("收到中止请求, 正在退出 ... ")
		// fmt.Println("退出。")
		os.Exit(0)
	}()
}

// 生成验证码
func GenerateCode(str string, timestamp int64, length int) string {
	// 将字符串和时间戳拼接起来
	data := fmt.Sprintf("%s%d", str, timestamp)

	// 计算MD5哈希值
	hasher := md5.New()
	hasher.Write([]byte(data))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// 随机获取6个字符
	codeStr := hash[:length]
	for i := 0; i < length; i++ {
		index := rand.Intn(len(hash))
		codeStr += string(hash[index])
	}

	// 取哈希值的前6位作为验证码
	code, _ := strconv.ParseInt(codeStr, 16, 64)
	var codeLen int64 = 1
	for i := 0; i < length; i++ {
		codeLen *= 10
	}
	code %= codeLen

	// 返回6位数字验证码
	fsf := fmt.Sprint("%0", length, "d")
	return fmt.Sprintf(fsf, code)
}

func (s Setting) SetLanguage(ishl bool, fl []string) int {
	if ishl {
		switch fl[0] {
		case "en", "1":
			return 1
		case "zh", "zh_cn", "zh-CN", "zhHans", "chs", "2":
			return 2
		case "zh_tw", "zh-TW", "zh_hk", "zh-HK", "zhHant", "cht", "3":
			return 3
		case "es", "4":
			return 4
		}
	}
	return s.Config.DefaultLocaleID
}

func (s Setting) LanguageMap(i18nStr string) map[string]string {
	i18ns := strings.Split(i18nStr, ",")
	temp := map[string]string{}
	for i := 0; i < len(s.Config.Language); i++ {
		if i >= len(i18ns) {
			break
		}
		temp[s.Config.Language[i]] = i18ns[i]
	}
	return temp
}

// 判断字符串中是否只有字母、数字和连字符
func IsLetterNumberHyphen(str string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	return re.MatchString(str)
}

func GenerateWhereTime(timeStr string) (string, []interface{}, int, error) {
	where := ""
	values := []interface{}{}
	tList := strings.Split(timeStr, ",")
	t1 := tList[0]
	t2 := tList[0]
	if len(tList) > 1 {
		t2 = tList[1]
	}

	timeType := ""
	_, err := time.Parse("2006-01-02", t1)
	if err == nil {
		timeType = "day"
	}
	_, err = time.Parse("2006-01", t1)
	if err == nil {
		timeType = "month"
	}
	switch timeType {
	case "day":
		where = "`time` BETWEEN ? AND ?"
		values = append(values, fmt.Sprintf("%s 00:00:00", t1))
		values = append(values, fmt.Sprintf("%s 23:59:59", t2))
	case "month":
		yearStr := t1[0:4]
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return where, values, 2041, fmt.Errorf("time")
		}
		month := t1[5:7]
		dayEnd := "31"
		switch month {
		case "02":
			if year%4 == 0 {
				dayEnd = "29"
			} else {
				dayEnd = "28"
			}
		case "04", "06", "09", "11":
			dayEnd = "30"
		case "01", "03", "05", "07", "08", "10", "12":
			dayEnd = "31"

		}
		where = "`time` BETWEEN ? AND ?"
		values = append(values, fmt.Sprintf("%s-01 00:00:00", t1))
		values = append(values, fmt.Sprintf("%s-%s 23:59:59", t2, dayEnd))
	}
	return where, values, -1, nil
}
