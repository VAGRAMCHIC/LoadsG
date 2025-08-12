package lib

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strings"
)

//func (m method) String() string {
//	return HttpMethods[m]
//}

func (m _httpMethod) IsValid() bool {
	switch m {
	case GET, POST, PUT, DELETE, CONNECT, PATCH, OPTIONS, HEAD:
		return true
	}
	return false
}

func parseRequestLine(url_path string) (method, host, path string, err error) {
	rawURL := url_path

	// Если url_path полный — разбираем через net/url_path
	u, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return "", "", "", parseErr
	}

	host = u.Host
	if host == "" {
		return "", "", "", fmt.Errorf("host не найден в request_line")
	}

	path = u.RequestURI()
	return method, host, path, nil
}

func CreateHttpHead(mt _httpMethod, url_path string, proto_version string, headers map[string]string) HttpHead {
	var request string
	var head HttpHead

	head.Method = mt
	if head.Method.IsValid() != true {
		panic("invalid_method")
	}
	head.URL = url_path
	head.ProtoVersion = proto_version
	head.Length = 0
	head.Headers = headers
	request += fmt.Sprintf("%s %s %s \n", head.Method, head.URL, head.ProtoVersion)

	return head
}

func BuildHttpRequest(head_data HttpHead, body string) (string, string) {
	request := fmt.Sprintf("%s %s HTTP/1.1\r\n", head_data.Method, head_data.URL)
	_, host, _, _ := parseRequestLine(head_data.URL)
	request += fmt.Sprintf("Host: %s\r\n", host)

	for k, v := range head_data.Headers {
		request += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	if body != "" {
		request += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	}

	request += "Connection: close\r\n\r\n"

	if body != "" {
		request += body
	}

	return request, host

}

func SendHttpRequest(request string, dialHost string) string {
	if !strings.Contains(dialHost, ":") {
		dialHost = dialHost + ":80"
	}

	conn, err := net.Dial("tcp", dialHost)
	if err != nil {
		fmt.Println("dial err:", err)
		return err.Error()
	}
	defer conn.Close()

	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("write err:", err)
		return err.Error()
	}

	reader := bufio.NewReader(conn)
	for {
		line, e := reader.ReadString('\n')
		if e != nil {
			break
		}
		fmt.Print(line)
	}

	return ""
}

func createHandleRespose() string {
	return ""
}
