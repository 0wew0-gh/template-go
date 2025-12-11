package tool

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// MARK: 创建文件夹
func CheckFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("%s:%v", path, err)
		}
	}
	return nil
}

// MARK: 保存文件
func SaveFileForReader(path string, name string, content io.Reader) error {
	notEmpty, content, err := ensureNotEmpty(content)
	if err != nil {
		return err
	}
	if !notEmpty {
		return fmt.Errorf("content is nil")
	}

	filePath := fmt.Sprintf("%s%s.eml", path, name)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err := f.ReadFrom(content); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	_ = f.Close()
	return nil
}

// MARK: 确保非空
func ensureNotEmpty(r io.Reader) (bool, io.Reader, error) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if err == io.EOF {
		return false, r, nil // 空
	}
	if err != nil {
		return false, r, err
	}

	// 把首字节拼回去，恢复原流
	newR := io.MultiReader(bytes.NewReader(buf[:n]), r)
	return true, newR, nil
}

// MARK: 打开文件,并返回文件内容字符串
func OpenFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MARK: 保存base64图片
func (s Setting) SaveBase64ImageToFile(id int, base64Str string, filePath string) error {
	// 解码base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}

	err = s.SaveStringToFile(id, data, filePath)
	return err
}

func (s Setting) SaveStringToFile(id int, bytes interface{}, filePath string) error {
	OssConfig, err := s.GetOssConfigForID(id)
	if err != nil {
		return err
	}

	client, err := oss.New(OssConfig.Endpoint, OssConfig.AccessKeyID, OssConfig.AccessKeySecret)
	if err != nil {
		// fmt.Println("oss.New err:", err)
		return err
	}
	ossbucket, err := client.Bucket(OssConfig.Bucket)
	if err != nil {
		// fmt.Println("client.Bucket err:", err)
		return err
	}

	is404 := false
	getData, err := ossbucket.GetObject(filePath)
	if err != nil {
		// fmt.Println("ossbucket.GetObject err:", err)
		if strings.Contains(err.Error(), "StatusCode=404") {
			is404 = true
		} else {
			return err
		}
	}
	if !is404 {

		isSame := false
		readData, err := io.ReadAll(getData)
		if err != nil {
			// fmt.Println("io.ReadAll err:", err)
			return err
		}
		switch byteType := bytes.(type) {
		case []byte:
			if condition := string(readData) == string(byteType); condition {
				isSame = true
			}
		case string:
			if condition := string(readData) == byteType; condition {
				isSame = true
			}
		}
		if isSame {
			return fmt.Errorf("文件已存在")
		}
	}

	str := ""
	switch byteType := bytes.(type) {
	case []byte:
		str = string(byteType)
	case string:
		str = byteType
	}
	err = ossbucket.PutObject(filePath, strings.NewReader(str))
	if err != nil {
		return err
	}
	return nil
}

func (s Setting) ReadStringFromFile(id int, filePath string) (string, error) {
	OssConfig, err := s.GetOssConfigForID(id)
	if err != nil {
		return "", err
	}
	client, err := oss.New(OssConfig.Endpoint, OssConfig.AccessKeyID, OssConfig.AccessKeySecret)
	if err != nil {
		// fmt.Println("oss.New err:", err)
		return "", err
	}
	ossbucket, err := client.Bucket(OssConfig.Bucket)
	if err != nil {
		// fmt.Println("client.Bucket err:", err)
		return "", err
	}

	getData, err := ossbucket.GetObject(filePath)
	if err != nil {
		// fmt.Println("ossbucket.GetObject err:", err)
		return "", err
	} else {
		// fmt.Println("getData:", getData)
		readData, err := io.ReadAll(getData)
		if err != nil {
			// fmt.Println("io.ReadAll err:", err)
			return "", err
		}
		return string(readData), nil
	}
}
