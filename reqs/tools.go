package reqs

import (
	"fmt"
	"io"
	"kss"
	"net/http"
	"net/url"
	"strings"
)

// MakeUrl 构造完整的请求URL
func MakeUrl(url string, params kss.S) string {
	if !strings.Contains(url, "?") {
		url += "?"
	} else {
		url += "&"
	}
	for key, value := range params {
		url += fmt.Sprintf("%s=%s&", key, value)
	}
	return url[:len(url)-1]
}

// FormDataEncode 对请求表单进行编码
func FormDataEncode(formData kss.S) []byte {
	encoded := make([]byte, 0)
	for key, value := range formData {
		encoded = append(encoded, []byte(fmt.Sprintf("%s=%s&", key, value))...)
	}
	encoded = encoded[:len(encoded)-1]
	return encoded
}

// GetProtocols 获取 http https 的字符串
func GetProtocols(protocol string) (string, string) {
	var p1, p2 string
	if strings.HasSuffix(protocol, "https") {
		p1 = "http" + strings.TrimPrefix(protocol, "https")
		p2 = protocol
	} else {
		p1 = protocol
		p2 = "https" + strings.TrimPrefix(protocol, "http")
	}
	return p1, p2
}

// ReqSetHeaders 为请求设置请求头
func ReqSetHeaders(req *http.Request, headers kss.S) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

// MakeProxyTransport 设置代理
func MakeProxyTransport(proxy string) (*http.Transport, error) {
	if proxy == "" {
		return &http.Transport{}, nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	return transport, nil
}

// Done 获取响应
func Done(client *http.Client, req *http.Request) (*Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{origin: resp, Content: body, Text: string(body)}, nil
}

// Do 发送请求（默认client）
func Do(req *http.Request, headers kss.S) (*Response, error) {
	ReqSetHeaders(req, headers)
	client := &http.Client{}
	return Done(client, req)
}

// ResetContentType 重新设置Content-Type
func ResetContentType(req *http.Request, headers kss.S, newValue string) {
	delete(headers, "Content-Type")
	req.Header.Set("Content-Type", newValue)
}
