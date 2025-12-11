package tool

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	level   int
	Info    *log.Logger
	Warning *log.Logger
	Debug   *log.Logger
	Err     *log.Logger
}

const (
	Error = iota
	Warning
	Debug
	Info
)

// 设置日志文件
func (l Logger) SetupLogFile(path string, oldFile *os.File) (Logger, *os.File) {
	var (
		file *os.File = nil
		err  error    = nil
	)
	// 获取当前日期
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	filePath := fmt.Sprintf("%s%s.log", path, currentDate)
	// 判断是否有日志文件夹
	if err := CheckFolder(path); err != nil {
		return l, nil
	}

	// 判断文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if oldFile != nil {
			// 不存在则说明是第二天
			l.Println(Error, "不存在则说明是第二天")
			oldFile.Close()
		}

		// 如果文件不存在，则创建文件
		l.Println(Error, "创建日志文件")
		file, err = os.Create(filePath)
		if err != nil {
			l.Println(Error, "创建日志文件失败:", err)
			return l, nil
		}
		l.Info = log.New(file, "[INFO] ", log.LstdFlags)
		l.Warning = log.New(file, "[Warning] ", log.LstdFlags)
		l.Debug = log.New(file, "[Debug] ", log.LstdFlags)
		l.Err = log.New(file, "[ERROR] ", log.LstdFlags)
		return l, file
	}

	if oldFile != nil {
		return l, oldFile
	}

	l.Println(Info, "打开日志文件")
	// 返回 *log.Logger 对象
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Println(Error, "打开日志文件失败:", err)
		return l, nil
	}
	l.Info = log.New(file, "[INFO] ", log.LstdFlags)
	l.Warning = log.New(file, "[Warning] ", log.LstdFlags)
	l.Debug = log.New(file, "[Debug] ", log.LstdFlags)
	l.Err = log.New(file, "[ERROR] ", log.LstdFlags)
	return l, file
}

func (l Logger) SetLevel(level int) Logger {
	l.level = level
	return l
}
func (l Logger) GetLevel(level int) int {
	return l.level
}

func (l Logger) Println(t int, v ...interface{}) (n int, err error) {
	switch t {
	case Info:
		if l.level < Info {
			return 0, nil
		}
		if l.Info == nil {
			// log.Println(v...)
			return 0, fmt.Errorf("info logger is not set")
		}
		l.Info.Println(v...)
		return 0, nil
	case Warning:
		if l.level < Warning {
			return 0, nil
		}
		if l.Warning == nil {
			// log.Println(v...)
			return 0, fmt.Errorf("warning logger is not set")
		}
		l.Warning.Println(v...)
		return 0, nil
	case Debug:
		if l.level < Debug {
			return 0, nil
		}
		if l.Debug == nil {
			// log.Println(v...)
			return 0, fmt.Errorf("debug logger is not set")
		}
		l.Debug.Println(v...)
		return 0, nil
	case Error:
		if l.level < Error {
			return 0, nil
		}
		if l.Err == nil {
			// log.Println(v...)
			return 0, fmt.Errorf("error logger is not set")
		}
		l.Err.Println(v...)
		return 0, nil
	}
	// log.Println(v...)
	return 0, fmt.Errorf("invalid log type")
}

func (l Logger) Printf(t int, format string, v ...interface{}) (n int, err error) {
	switch t {
	case Info:
		if l.level < Info {
			return 0, nil
		}
		if l.Info == nil {
			// log.Printf(format, v...)
			return 0, fmt.Errorf("info logger is not set")
		}
		l.Info.Printf(format, v...)
		return 0, nil
	case Warning:
		if l.level < Warning {
			return 0, nil
		}
		if l.Warning == nil {
			// log.Printf(format, v...)
			return 0, fmt.Errorf("warning logger is not set")
		}
		l.Warning.Printf(format, v...)
		return 0, nil
	case Debug:
		if l.level < Debug {
			return 0, nil
		}
		if l.Debug == nil {
			// log.Printf(format, v...)
			return 0, fmt.Errorf("debug logger is not set")
		}
		l.Debug.Printf(format, v...)
		return 0, nil
	case Error:
		if l.level < Error {
			return 0, nil
		}
		if l.Err == nil {
			// log.Printf(format, v...)
			return 0, fmt.Errorf("error logger is not set")
		}
		l.Err.Printf(format, v...)
		return 0, nil

	}
	// log.Printf(format, v...)
	return 0, fmt.Errorf("invalid log type")
}
