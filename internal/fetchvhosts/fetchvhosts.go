//Package fetchvhosts logs in to server via SSH, retrieves current HTTP server configuration containing vhosts
//Retrieved configration is then processed and parsed into string slice
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

//FetchVhosts ...
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
	separatedStr := strings.FieldsFunc(strippedStr, func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	for i := 0; i < len(separatedStr); i++ {
		separatedStr[i] = strings.TrimSpace(separatedStr[i])
		if len(separatedStr[i]) > 0 {
			if separatedStr[i][0] == '#' || separatedStr[i][0] == '_' {
				separatedStr = append(separatedStr[:i], separatedStr[i+1:]...)
			}
		}
		separatedStr[i] = strings.Trim(separatedStr[i], ";")
	}

	uniqStr := removeDuplicates(separatedStr)

	cleanedStr := strings.Join(uniqStr, "\n")
	return strings.FieldsFunc(cleanedStr, func(r rune) bool {
		return r == ' ' || r == '\n'
	})
}

func getHttpdVhosts(session *ssh.Session) []string {
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run("sudo httpd -S")
	return strings.Split(stdoutBuf.String(), "\n")
}

func removeDuplicates(str []string) []string {
	dict := make(map[string]int)
	for i := range str {
		dict[str[i]]++
	}
	result := make([]string, len(dict))
	for k := range dict {
		result = append(result, k)
	}
	return result
}
