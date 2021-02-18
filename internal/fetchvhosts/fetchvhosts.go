package fetchvhosts

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

var serverMap = map[string]func(*ssh.Session) []string{
	"nginx":  getNginxVhosts,
	"apache": getHttpdVhosts,
}

func FetchVhosts(host, httpServer string, config *ssh.ClientConfig) []string {
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	vhosts := make([]string, 0)
	if _, ok := serverMap[httpServer]; ok {
		vhosts = serverMap[httpServer](session)
	}

	return vhosts
}

func getNginxVhosts(session *ssh.Session) []string {
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run("sudo -n nginx -T | grep server_name")

	strippedStr := string(bytes.ReplaceAll(stdoutBuf.Bytes(), []byte("server_name"), []byte("")))
	separatedStr := strings.Split(strippedStr, "\n")

	for i := 0; i < len(separatedStr); i++ {
		if len(separatedStr[i]) > 0 {
			if separatedStr[i][0] == '#' || separatedStr[i][0] == '_' {
				separatedStr = append(separatedStr[:i], separatedStr[i+1:]...)
			}
		}
	}

	cleanedStr := strings.Join(separatedStr, "\n")
	return strings.FieldsFunc(cleanedStr, func(r rune) bool {
		return r == ' ' || r == '\n'
	})
	// return strings.Split(stdoutBuf.String(), "\n")
}

func getHttpdVhosts(session *ssh.Session) []string {
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run("sudo httpd -S")
	return strings.Split(stdoutBuf.String(), "\n")
}
