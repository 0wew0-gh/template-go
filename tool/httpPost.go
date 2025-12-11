package tool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// http请求,根据设置使用代理
func (s Setting) HTTPPost(urlPath string, key []string, value []string) (map[string]interface{}, error) {
	if len(key) != len(value) {
		return nil, fmt.Errorf("key and value length not equal")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	proxyURL, err := url.Parse(s.Config.Proxy)
	if err != nil {
		proxyURL = nil
	}
	if proxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	formData := url.Values{}
	for i := 0; i < len(key); i++ {
		formData.Set(key[i], value[i])
	}
	resp, err := client.PostForm(urlPath, formData)
	if err != nil {
		return nil, err
	}
	body, err := UnMarshalBody(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}

func UnMarshalBody(body io.ReadCloser) (map[string]interface{}, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %v", err)
	}
	return data, nil
}
