package tool

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// ==========
//
//	发送邮件
//	to	string		目标邮件地址
//	subject	string		主题
//	body	string		内容
//	errC	chan error	错误
func (s Setting) SendMail(to string, subject string, body string, errC chan error) {
	auth := smtp.PlainAuth("", s.Config.Mail.From, s.Config.Mail.PW, s.Config.Mail.Host)

	// fmt.Println(s.Config.Mail.From, s.Config.Mail.FormName, s.Config.Mail.Host, subject)
	var content_type string
	if s.Config.Mail.Type == "html" {
		content_type = "text/" + s.Config.Mail.Type + "; charset=UTF-8"
	} else {
		content_type = "text/plain" + "; charset=UTF-8"
	}

	curtime := time.Now().In(s.Config.TimeZone.Location).Format(s.Config.TimeZone.Format)

	header := make(map[string]string)
	header["From"] = s.Config.Mail.FormName + "<" + s.Config.Mail.From + ">"
	header["To"] = to
	header["Date"] = curtime
	header["Subject"] = subject
	header["Content-Type"] = content_type
	// msg := []byte("To: " + to[0] + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	msg += "\r\n" + body
	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%s", s.Config.Mail.Host, s.Config.Mail.Port),
		auth,
		s.Config.Mail.From,
		[]string{to},
		[]byte(msg),
	)
	if err != nil {
		errC <- err
		return
	}
	errC <- err
}

// return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		// fmt.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		// fmt.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				// fmt.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
