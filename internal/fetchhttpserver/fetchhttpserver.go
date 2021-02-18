package fetchhttpserver

import (
	"net/http"
	"strings"
)

func evaluateHTTPServer(header []string) string {
	httpServers := []string{
		"apache",
		"httpd",
		"nginx",
		"litespeed",
	}

	for i := range header {
		for j := range httpServers {
			if strings.Contains(strings.ToLower(header[i]), httpServers[j]) {
				return httpServers[j]
			}
		}
	}
	return ""
}

func checkHTTPServer(host string) ([]string, error) {
	url := strings.Split(host, ":")[0]

	res, err := http.Head("http://" + url)
	if err != nil {
		return nil, err
	}

	return res.Header["Server"], nil
}

func FetchHTTPServer(host string) (string, error) {
	httpHeader, err := checkHTTPServer(host)
	if err != nil {
		return "", err
	}
	return evaluateHTTPServer(httpHeader), nil
}
